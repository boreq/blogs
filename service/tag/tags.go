package tag

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/sqlx"
)

var log = logging.New("service/tags")

func New(db *sqlx.DB) *TagService {
	rv := &TagService{
		db: db,
	}
	return rv
}

type TagService struct {
	db *sqlx.DB
}

func (b *TagService) GetForPost(postId uint) ([]database.Tag, error) {
	log.Debug("GetForPost", "postId", postId)
	query := `SELECT tag.*
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id = tag.id
		JOIN post ON post.id = post_to_tag.post_id
		WHERE post.id=$1`

	tags := make([]database.Tag, 0)
	if err := database.DB.Select(&tags, query, postId); err != nil {
		return nil, err
	}
	return tags, nil
}
