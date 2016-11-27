package database

import (
	"fmt"
	"github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/utils"
	"time"
)

type Blog struct {
	ID uint `gorm:"primary_key"`

	InternalID uint   `gorm:"not null"`
	Title      string `gorm:"not null"`
	Categories []Category
}

func (blog *Blog) GetUrl() string {
	loader, ok := blogs.Blogs[blog.InternalID]
	if ok {
		return loader.GetUrl()
	}
	return ""
}

func (blog *Blog) GetAbsoluteUrl() string {
	return fmt.Sprintf("/blog/%d/%s", blog.ID, utils.Slugify(blog.Title))
}

type Category struct {
	ID     uint `gorm:"primary_key"`
	BlogID uint

	Name  string `gorm:"not null"`
	Posts []Post
}

type Post struct {
	ID         uint `gorm:"primary_key"`
	CategoryID uint

	InternalID string `gorm:"not null"`
	Title      string `gorm:"not null"`
	Date       time.Time
	Tags       []Tag `gorm:"many2many:post_tags;"`
}

func (post Post) GetUrl() string {
	category := &Category{}
	blog := &Blog{}
	DB.Model(post).Related(&category)
	DB.Model(category).Related(&blog)
	loader, ok := blogs.Blogs[blog.InternalID]
	if ok {
		return loader.GetPostUrl(post.InternalID)
	}
	return ""
}

func (post *Post) GetCategory() *Category {
	category := &Category{}
	DB.Model(&post).Related(category)
	return category
}

type Tag struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null"`
}
