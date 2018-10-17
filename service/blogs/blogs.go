package blogs

import (
	"database/sql"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/dto"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/sqlx"
	"github.com/pkg/errors"
)

var log = logging.New("service/blogs")

func New(db *sqlx.DB) *BlogsService {
	rv := &BlogsService{
		db: db,
	}
	return rv
}

type BlogsService struct {
	db *sqlx.DB
}

const (
	SortTitle       ListSort = "blog.title"
	SortSubscribers ListSort = "blog.subscriptions"
	SortLastPost    ListSort = "last_post"
)

type ListSort string

type blogResult struct {
	database.Blog
	LastPost   *dto.ScannableTime
	Subscribed sql.NullInt64
}

type ListOut struct {
	Page  dto.PageOut   `json:"page"`
	Blogs []dto.BlogOut `json:"blogs"`
}

func (b *BlogsService) List(page dto.Page, sort ListSort, reverse bool, userId *uint) (ListOut, error) {
	var amount uint
	if err := b.db.Get(&amount, "SELECT COUNT(*) AS amount FROM blog"); err != nil {
		return ListOut{}, errors.Wrap(err, "could not count the blogs")
	}

	limit, offset := limitOffset(page)

	var blogs []blogResult
	if err := b.db.Select(&blogs,
		`SELECT blog.*, MAX(post.date) AS last_post, subscription.id AS subscribed
		FROM blog
		LEFT JOIN category ON category.blog_id=blog.id
		LEFT JOIN post ON post.category_id=category.id
		LEFT JOIN subscription ON subscription.blog_id=blog.id AND subscription.user_id=$1
		GROUP BY blog.id, subscription.id
		ORDER BY `+string(sort)+` `+order(reverse)+`
		LIMIT $2 OFFSET $3`, userId, limit, offset); err != nil {
		return ListOut{}, errors.Wrap(err, "could not get the blogs")
	}

	blogsOut, err := toBlogsOut(blogs)
	if err != nil {
		return ListOut{}, errors.Wrap(err, "could not convert to blogs out")
	}
	out := ListOut{
		Page: dto.PageOut{
			Page:     page,
			AllItems: int(amount),
		},
		Blogs: blogsOut,
	}
	return out, nil
}

func toBlogsOut(blogResults []blogResult) ([]dto.BlogOut, error) {
	blogsOut := make([]dto.BlogOut, 0)
	for _, blogResult := range blogResults {
		subscribed := blogResult.Subscribed.Valid && blogResult.Subscribed.Int64 > 0
		url, err := blogResult.Blog.GetUrl()
		if err != nil {
			return nil, errors.Wrapf(err, "could not get the url for blog %+v", blogResult.Blog)
		}
		cleanUrl, err := blogResult.Blog.GetCleanUrl()
		if err != nil {
			return nil, errors.Wrapf(err, "could not get the clean url for blog %+v", blogResult.Blog)
		}
		blogOut := dto.BlogOut{
			Blog:       blogResult.Blog,
			LastPost:   blogResult.LastPost,
			Subscribed: subscribed,
			Url:        url,
			CleanUrl:   cleanUrl,
		}
		blogsOut = append(blogsOut, blogOut)
	}
	return blogsOut, nil
}

func order(reverse bool) string {
	if reverse {
		return "DESC"
	} else {
		return "ASC"
	}
}

func limitOffset(page dto.Page) (int, int) {
	offset := page.PerPage * (page.Page - 1)
	limit := page.PerPage
	return limit, offset

}
