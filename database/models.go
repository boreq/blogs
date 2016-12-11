package database

import (
	"fmt"
	"github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/utils"
	"time"
)

type User struct {
	ID uint

	Username string
	Password string
}

type UserSession struct {
	UserID uint

	Key      string
	LastSeen time.Time
}

type Blog struct {
	ID uint

	InternalID    uint
	Title         string
	Subscriptions int
}

func (blog Blog) GetUrl() string {
	loader, ok := blogs.Blogs[blog.InternalID]
	if ok {
		return loader.GetUrl()
	}
	return ""
}

func (blog Blog) GetAbsoluteUrl() string {
	return fmt.Sprintf("/blog/%d/%s", blog.ID, utils.Slugify(blog.Title))
}

type Category struct {
	ID     uint
	BlogID uint

	Name string
}

type Post struct {
	ID         uint
	CategoryID uint

	InternalID string
	Title      string
	Summary    string
	Date       time.Time
	Stars      int
}

func (post Post) GetISO8601Date() string {
	return post.Date.Format(time.RFC3339)
}

type Tag struct {
	ID   uint
	Name string
}

type Update struct {
	ID     uint
	BlogID uint

	Started   time.Time
	Ended     time.Time
	Succeeded bool
	Data      string
}

type Subscription struct {
	ID     uint
	BlogID uint
	UserID uint
	Date   time.Time
}

type Star struct {
	ID     uint
	PostID uint
	UserID uint
	Date   time.Time
}
