package updater

import (
	"github.com/boreq/blogs/blogs/loaders"
	"github.com/boreq/blogs/database"
	"time"
)

func newBlogUpdater(blogDatabaseID uint, loader loaders.Blog) *blogUpdater {
	rv := &blogUpdater{
		blogDatabaseID: blogDatabaseID,
		loader:         loader,
		stats:          statistics{},
	}
	return rv
}

type statistics struct {
	NumPosts    int
	Errors      []error
	ElapsedTime time.Duration
}

type blogUpdater struct {
	blogDatabaseID uint
	loader         loaders.Blog
	stats          statistics
	blog           *database.Blog
}

// Run performs a full update of the blog.
func (u *blogUpdater) Run() {
	start := time.Now()

	u.blog = u.updateBlogDatabaseEntry()
	u.updatePosts()

	u.stats.ElapsedTime = time.Since(start)
}

// updateBlogDatabaseEntry ensures that there is a proper entry in the blogs
// table and returns it.
func (u *blogUpdater) updateBlogDatabaseEntry() *database.Blog {
	blog := &database.Blog{}
	database.DB.FirstOrInit(&blog, &database.Blog{InternalID: u.blogDatabaseID})

	// Update the title
	title, err := u.loader.LoadTitle()
	if err != nil {
		u.stats.Errors = append(u.stats.Errors, err)
	} else {
		blog.Title = title
	}

	// Save or create the record
	if database.DB.NewRecord(blog) {
		database.DB.Create(blog)
	} else {
		database.DB.Save(blog)
	}

	return blog
}

// getOrCreateCategory returns an existing category or creates it.
func (u *blogUpdater) getOrCreateCategory(name string) *database.Category {
	category := &database.Category{}
	database.DB.FirstOrCreate(&category, &database.Category{BlogID: u.blog.ID, Name: name})
	return category
}

// getOrCreateTag returns an existing tag or creates it.
func (u *blogUpdater) getOrCreateTag(name string) *database.Tag {
	tag := &database.Tag{}
	database.DB.FirstOrCreate(&tag, &database.Tag{Name: name})
	return tag
}

// updatePosts downloads and updates all posts made on the blog.
func (u *blogUpdater) updatePosts() {
	posts, errors := u.loader.LoadPosts()
	for {
		select {
		case err, ok := <-errors:
			if ok {
				u.stats.Errors = append(u.stats.Errors, err)
			} else {
				errors = nil
			}
		case post, ok := <-posts:
			if ok {
				u.stats.NumPosts += 1
				u.handlePost(post)
			} else {
				posts = nil
			}
		}
		if errors == nil && posts == nil {
			break
		}
	}
}

// handlePost updates a single post received from the loader.
func (u *blogUpdater) handlePost(loadedPost loaders.Post) {
	post := &database.Post{}
	database.DB.Preload("Tags").FirstOrInit(&post, &database.Post{InternalID: loadedPost.Id})

	altered := false
	if post.Date != loadedPost.Date {
		post.Date = loadedPost.Date
		altered = true
	}
	if post.Title != loadedPost.Title {
		post.Title = loadedPost.Title
		altered = true
	}
	if post.Summary != loadedPost.Summary {
		post.Summary = loadedPost.Summary
		altered = true
	}
	category := u.getOrCreateCategory(loadedPost.Category)
	if post.CategoryID != category.ID {
		post.CategoryID = category.ID
		altered = true
	}

	if database.DB.NewRecord(post) {
		database.DB.Create(post)
	} else if altered {
		database.DB.Save(post)
	}

	u.addNewTags(post, loadedPost.Tags, post.Tags)
	u.removeOldTags(post, loadedPost.Tags, post.Tags)
}

// addNewTags ensures that all tags are present in the database.
func (u *blogUpdater) addNewTags(post *database.Post, loadedTags []string, tags []database.Tag) {
	for _, loadedTagName := range loadedTags {
		loadedTag := u.getOrCreateTag(loadedTagName)
		if tag := u.findTag(loadedTag, tags); tag == nil {
			database.DB.Model(&post).Association("Tags").Append(loadedTag)
		}
	}
}

func (u *blogUpdater) findTag(loadedTag *database.Tag, tags []database.Tag) *database.Tag {
	for _, tag := range tags {
		if tag.ID == loadedTag.ID {
			return &tag
		}
	}
	return nil
}

// removeOldTags ensures that all old tags are removed from the database.
func (u *blogUpdater) removeOldTags(post *database.Post, loadedTags []string, tags []database.Tag) {
	for _, tag := range tags {
		if !u.findLoadedTag(tag.Name, loadedTags) {
			database.DB.Model(&post).Association("Tags").Delete(tag)
		}
	}
}

func (u *blogUpdater) findLoadedTag(name string, tags []string) bool {
	for _, tag := range tags {
		if tag == name {
			return true
		}
	}
	return false
}
