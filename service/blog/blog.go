package blog

import (
	"database/sql"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/dto"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/sqlx"
	"github.com/pkg/errors"
	"time"
)

var log = logging.New("service/blog")

func New(db *sqlx.DB) *BlogService {
	rv := &BlogService{
		db: db,
	}
	return rv
}

type BlogService struct {
	db *sqlx.DB
}

func (b *BlogService) Subscribe(blogId uint, userId uint) error {
	if _, err := b.db.Exec(`
		INSERT INTO subscription(blog_id, user_id, date)
		SELECT $1, $2, $3
		WHERE NOT EXISTS(
			SELECT 1
			FROM subscription
			WHERE blog_id=$1 AND user_id=$2)`,
		blogId, userId, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (b *BlogService) Unsubscribe(blogId uint, userId uint) error {
	if _, err := b.db.Exec(`
		DELETE FROM subscription
		WHERE blog_id=$1 AND user_id=$2`,
		blogId, userId); err != nil {
		return err
	}
	return nil
}

type blogResult struct {
	database.Blog
	Subscribed sql.NullInt64
}

func (b *BlogService) Get(blogId uint, userId *uint) (dto.BlogOut, error) {
	query := `SELECT blog.*, subscription.id AS subscribed
	FROM blog
	LEFT JOIN subscription ON subscription.blog_id=blog.id AND subscription.user_id=$1
	WHERE blog.id=$2`

	var blog blogResult
	err := database.DB.Get(&blog, query, userId, blogId)
	if err != nil {
		return dto.BlogOut{}, errors.Wrap(err, "could not get the blog from database")
	}
	return toBlogOut(blog)
}

func toBlogOut(blogResult blogResult) (dto.BlogOut, error) {
	subscribed := blogResult.Subscribed.Valid && blogResult.Subscribed.Int64 > 0
	url, err := blogResult.Blog.GetUrl()
	if err != nil {
		return dto.BlogOut{}, errors.Wrapf(err, "could not get the url for blog %+v", blogResult.Blog)
	}
	cleanUrl, err := blogResult.Blog.GetCleanUrl()
	if err != nil {
		return dto.BlogOut{}, errors.Wrapf(err, "could not get the clean url for blog %+v", blogResult.Blog)
	}
	blogOut := dto.BlogOut{
		Blog:       blogResult.Blog,
		Subscribed: subscribed,
		Url:        url,
		CleanUrl:   cleanUrl,
	}
	return blogOut, nil
}

func (b *BlogService) GetCategories(blogId uint) ([]database.Category, error) {
	query := `SELECT category.*
		FROM category
		JOIN blog ON blog.id=category.blog_id
		WHERE blog.id=$1`
	var categories []database.Category
	if err := database.DB.Select(&categories, query, blogId); err != nil {
		return nil, errors.Wrap(err, "could not get categories")
	}
	return categories, nil
}

func (b *BlogService) GetTags(blogId uint) ([]dto.TagOut, error) {
	query := `SELECT tag.*, COUNT(tag.id) as posts
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id = tag.id
		JOIN post ON post.id = post_to_tag.post_id
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		WHERE blog.id=$1
		GROUP BY tag.id`
	var tags []dto.TagOut
	if err := database.DB.Select(&tags, query, blogId); err != nil {
		return nil, errors.Wrap(err, "could not get tags")
	}
	return tags, nil
}
