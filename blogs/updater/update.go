package updater

import (
	"github.com/boreq/blogs/blogs"
)

// Update performs an update of all defined blogs.
func Update() error {
	for blogDatabaseID, blog := range blogs.Blogs {
		updater := newBlogUpdater(blogDatabaseID, blog)
		updater.Run()
	}
	return nil
}
