package blog

import (
	"github.com/boreq/blogs/logging"
	"github.com/boreq/sqlx"
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
