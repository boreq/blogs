package loaders

import (
	"time"
)

type Blog interface {
	// GetUrl returns an url pointing to the home page of the blog. The url
	// should not include the scheme.
	GetUrl() string

	// GetPostUrl returns an url pointing to the specific post on the blog.
	// The url should not include the scheme.
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
	Tags     []string
}
