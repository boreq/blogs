package database

import (
	"fmt"
	"github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model

	Username string `gorm:"size:255"`
	Password string `gorm:"size:255"`
	Sessions []UserSession
}

type UserSession struct {
	gorm.Model
	UserID uint `gorm:"not null"`

	SessionKey string `gorm:"size:512"`
	LastSeen   time.Time
}

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

	InternalID string `gorm:"not null;size:1000;"`
	Title      string `gorm:"not null;size:1000;"`
	Summary    string `gorm:"not null;size:3000;"`
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

func (post Post) GetISO8601Date() string {
	return post.Date.Format(time.RFC3339)
}

type Tag struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null"`
}
