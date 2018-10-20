package tag

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/dto"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/posts"
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

func (b *TagService) AddTags(listOut posts.ListOut) (ListOutWithTags, error) {
	postsOut, err := b.AddTagsToPosts(listOut.Posts)
	if err != nil {
		return ListOutWithTags{}, err
	}
	rv := ListOutWithTags{
		Page:  listOut.Page,
		Posts: postsOut,
	}
	return rv, nil
}

func (b *TagService) AddTagsToPosts(posts []dto.PostOut) ([]PostOutWithTags, error) {
	postsOut := make([]PostOutWithTags, 0)
	for _, postOut := range posts {
		tags, err := b.GetForPost(postOut.Post.ID)
		if err != nil {
			return nil, err
		}
		postsOut = append(postsOut, PostOutWithTags{postOut, tags})
	}
	return postsOut, nil
}

type PostOutWithTags struct {
	dto.PostOut
	Tags []database.Tag `json:"tags"`
}

type ListOutWithTags struct {
	Page  dto.PageOut       `json:"page"`
	Posts []PostOutWithTags `json:"posts"`
}
