package post

import (
	"github.com/boreq/blogs/logging"
	"github.com/boreq/sqlx"
	"time"
)

var log = logging.New("service/post")

func New(db *sqlx.DB) *PostService {
	rv := &PostService{
		db: db,
	}
	return rv
}

type PostService struct {
	db *sqlx.DB
}

func (p *PostService) Star(postId uint, userId uint) error {
	log.Debug("starring", "postId", postId, "userId", userId)
	if _, err := p.db.Exec(`
		INSERT INTO star(post_id, user_id, date)
		SELECT $1, $2, $3
		WHERE NOT EXISTS(
			SELECT 1
			FROM star
			WHERE post_id=$1 AND user_id=$2)`,
		postId, userId, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (p *PostService) Unstar(postId uint, userId uint) error {
	log.Debug("unstarring", "postId", postId, "userId", userId)
	if _, err := p.db.Exec(`
		DELETE FROM star
		WHERE post_id=$1 AND user_id=$2`,
		postId, userId); err != nil {
		return err
	}
	return nil
}
