package core

import (
	"database/sql"
	"errors"
	blgs "github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/utils"
	"time"
)

type postCategoryBlog struct {
	database.Post
	database.Category
	database.Blog
}

func (p postCategoryBlog) GetPostUrl() string {
	loader, ok := blgs.Blogs[p.Blog.InternalID]
	if ok {
		return loader.GetPostUrl(p.Post.InternalID)
	}
	return ""
}

func (p postCategoryBlog) GetPostTags() []database.Tag {
	var tags []database.Tag
	err := database.DB.Select(&tags,
		`SELECT tag.*
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id = tag.id
		JOIN post ON post.id = post_to_tag.post_id
		WHERE post.id=$1
		ORDER BY tag.name DESC`, p.Post.ID)
	if err != nil {
		panic(err)
	}
	return tags
}

type postsResult struct {
	postCategoryBlog
	Starred sql.NullInt64
}

type tagResult struct {
	database.Tag
	Count uint
}

type scannableTime struct {
	time.Time
}

func (t *scannableTime) Scan(src interface{}) error {
	updated, ok := src.([]uint8)
	if !ok {
		return errors.New("Invalid type, this is not []uint8")
	}
	tmp, err := time.Parse("2006-01-02 15:04:05-07:00", string(updated))
	if err != nil {
		return err
	}
	t.Time = tmp
	return nil
}

func (t scannableTime) String() string {
	return utils.ISO8601(t.Time)
}

type blogResult struct {
	database.Blog
	Updated    scannableTime
	Subscribed sql.NullInt64
}

type updateResult struct {
	database.Update
	database.Blog
}
