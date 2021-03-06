// Package loaders implements interfaces and data structures used for
// downloading posts and other data from the blogs.
package loaders

import (
	"time"
)

type Blog interface {
	// GetUrl returns an url pointing to the home page of the blog.
	GetUrl() string

	// GetCleanUrl returns an url which looks nice presented to the user.
	GetCleanUrl() string

	// GetPostUrl returns an url pointing to the specific post on the blog.
	GetPostUrl(internalID string) string

	// LoadTitle downloads a title of the blog.
	LoadTitle() (string, error)

	// LoadPosts downloads all posts made on the blog and sends them on the
	// returned channel. The second channel is used for sending errors and
	// both posts and errors can be sent at the same time. Both channels
	// will be closed once all posts and errors are read.
	LoadPosts() (<-chan Post, <-chan error)
}

type Post struct {
	Id       string
	Date     time.Time
	Category string
	Title    string
	Summary  string
	Tags     []string
}
