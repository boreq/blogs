package posts

import (
	"database/sql"
	"github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/dto"
	"github.com/boreq/blogs/logging"
	sqlutils "github.com/boreq/blogs/utils/sql"
	"github.com/boreq/sqlx"
	"github.com/pkg/errors"
)

var log = logging.New("service/posts")

func New(db *sqlx.DB) *PostsService {
	rv := &PostsService{
		db: db,
	}
	return rv
}

type PostsService struct {
	db *sqlx.DB
}

const (
	SortDate  ListSort = "post.date"
	SortStars ListSort = "post.stars"
	SortTitle ListSort = "post.title"
)

type ListSort string

type postCategoryBlog struct {
	database.Post
	database.Category
	database.Blog
}

type postResult struct {
	postCategoryBlog
	Starred sql.NullInt64
}

type ListOut struct {
	Page  dto.PageOut   `json:"page"`
	Posts []dto.PostOut `json:"posts"`
}

func (p *PostsService) List(page dto.Page, sort ListSort, reverse bool, userId *uint) (ListOut, error) {
	var amount uint
	if err := p.db.Get(&amount, "SELECT COUNT(*) AS amount FROM post"); err != nil {
		return ListOut{}, errors.Wrap(err, "could not count the posts")
	}

	limit, offset := sqlutils.LimitOffset(page)

	query := `SELECT post.*, category.*, blog.*, star.id AS starred
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		LEFT JOIN star ON star.post_id=post.id AND star.user_id=$1
		GROUP BY post.id, category.id, blog.id, star.id
		ORDER BY ` + string(sort) + ` ` + sqlutils.Order(reverse) + `
		LIMIT $2 OFFSET $3`

	var posts []postResult
	if err := p.db.Select(&posts, query, userId, limit, offset); err != nil {
		return ListOut{}, errors.Wrap(err, "could not get the posts")
	}

	postsOut, err := toPostsOut(posts)
	if err != nil {
		return ListOut{}, errors.Wrap(err, "could not convert to posts out")
	}
	out := ListOut{
		Page: dto.PageOut{
			Page:     page,
			AllItems: int(amount),
		},
		Posts: postsOut,
	}
	return out, nil
}

func toPostsOut(postResults []postResult) ([]dto.PostOut, error) {
	postsOut := make([]dto.PostOut, 0)
	for _, postResult := range postResults {
		starred := postResult.Starred.Valid && postResult.Starred.Int64 > 0
		url, err := getPostUrl(postResult.Blog, postResult.Post)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get the url for post %+v", postResult.Post)
		}
		postOut := dto.PostOut{
			Post:     postResult.Post,
			Category: postResult.Category,
			Blog:     postResult.Blog,
			Url:      url,
			Starred:  starred,
		}
		postsOut = append(postsOut, postOut)
	}
	return postsOut, nil
}

// GetUrl returns the address of the post.
func getPostUrl(blog database.Blog, post database.Post) (string, error) {
	loader, ok := blogs.Blogs[blog.InternalID]
	if ok {
		return loader.GetPostUrl(post.InternalID), nil
	}
	return "", errors.New("loader could not be found for this blog")
}
