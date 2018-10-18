package database

import (
	"github.com/boreq/blogs/blogs"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID uint `json:"id"`

	Username string `json:"username"`
	Password string `json:"-"`
}

type UserSession struct {
	ID     uint
	UserID uint

	Key      string
	LastSeen time.Time
}

type Blog struct {
	ID uint `json:"id"`

	InternalID    uint   `json:"-"`
	Title         string `json:"title"`
	Subscriptions int    `json:"subscriptions"`
}

// GetUrl returns the address of the blog.
func (blog Blog) GetUrl() (string, error) {
	loader, ok := blogs.Blogs[blog.InternalID]
	if ok {
		return loader.GetUrl(), nil
	}
	return "", errors.New("loader could not be found for this blog")
}

// GetCleanUrl returns the address of the blog which looks nice presented to
// the user.
func (blog Blog) GetCleanUrl() (string, error) {
	loader, ok := blogs.Blogs[blog.InternalID]
	if ok {
		return loader.GetCleanUrl(), nil
	}
	return "", errors.New("loader could not be found for this blog")
}

type Category struct {
	ID     uint `json:"id"`
	BlogID uint `json:"-"`

	Name string `json:"name"`
}

type Post struct {
	ID         uint `json:"id"`
	CategoryID uint `json:"-"`

	InternalID string    `json:"-"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Date       time.Time `json:"date"`
	Stars      int       `json:"stars"`
}

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
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
