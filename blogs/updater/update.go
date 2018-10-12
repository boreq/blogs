package updater

import (
	"errors"
	"fmt"
	"github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/blogs/loaders"
	"github.com/boreq/blogs/logging"
	"time"
)

var log = logging.New("updater")

// Update performs an update of all defined blogs.
func Update() error {
	for internalID, loader := range blogs.Blogs {
		updater := newBlogUpdater(internalID, loader)
		updater.Run()
	}
	return nil
}

// UpdateSpecific performs an update of the specific blog.
func UpdateSpecific(internalID uint) error {
	loader, ok := blogs.Blogs[internalID]
	if !ok {
		return errors.New("No such blog")
	}
	updater := newBlogUpdater(internalID, loader)
	updater.Run()
	return nil
}

// TestLoader doesn't alter the database in any way and only prints the
// downloaded data.
func TestLoader(internalID uint) error {
	loader, ok := blogs.Blogs[internalID]
	if !ok {
		return errors.New("No such blog")
	}

	numPosts := 0
	numErrors := 0

	start := time.Now()
	title, err := loader.LoadTitle()
	if err != nil {
		log.Error("Error", "err", err)
		numErrors++
	}
	posts, errors := loader.LoadPosts()
	for {
		select {
		case err, ok := <-errors:
			if ok {
				log.Error("Error", "err", err)
				numErrors++
			} else {
				errors = nil
			}
		case post, ok := <-posts:
			if ok {
				log.Info("Post", "post", postToString(post))
				numPosts++
			} else {
				posts = nil
			}
		}
		if errors == nil && posts == nil {
			break
		}
	}
	elapsed := time.Since(start)

	log.Info("Title", "title", title)
	log.Info("Errors", "errors", numErrors)
	log.Info("Posts", "posts", numPosts)
	log.Info("Elapsed", "elapsed", elapsed)
	return nil
}

func postToString(post loaders.Post) string {
	return fmt.Sprintf("Id: %s, Title: %s, Date: %s, Category: %s, Tags: %s, Summary: %s",
		post.Id, post.Title, post.Date, post.Category, post.Tags, post.Summary)
}
