package blogs

import (
	"github.com/boreq/blogs/logging"
)

var log = logging.GetLogger("blogs")

// Update performs an update of all defined blogs.
func Update() error {
	for blogDatabaseID, blog := range Blogs {
		updater := newBlogUpdater(blogDatabaseID, blog)
		updater.Run()
	}
	return nil
}
