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
	}
	return rv
}

type blogUpdater struct {
	blogDatabaseID uint
	loader         loaders.Blog
	stats          update
}

// Run performs a full update of the blog.
func (u *blogUpdater) Run() {
	u.stats.Started = time.Now()

	blog, err := getBlog(u.blogDatabaseID)
	if err == nil {
		u.updateBlogDatabaseEntry(blog)
		u.updatePosts(blog)

		u.stats.Ended = time.Now()
		u.stats.BlogID = blog.ID
		if err := saveStatistics(u.stats); err != nil {
			u.log("update statistics weren't saved: %s", err)
		}
	}
}

// updateBlogDatabaseEntry updates the blog title. If this function returns an
// error the update process can continue.
func (u *blogUpdater) updateBlogDatabaseEntry(blog database.Blog) {
	// Download
	title, err := u.loader.LoadTitle()
	if err != nil {
		u.log("blog title wasn't downloaded - loader error: %s", err)
		return
	}
	u.stats.TitleDownloaded = true

	// Update
	if blog.Title == title {
		u.stats.TitleCorrect = true
	} else {
		blog.Title = title

		if _, err := database.DB.NamedExec("UPDATE blog SET title=:title WHERE internal_id=:internal_id", blog); err != nil {
			u.log("blog title wasn't saved - database error: %s", err)
		}
		u.stats.TitleUpdated = true
	}
}

// updatePosts downloads and updates all posts made on the blog.
func (u *blogUpdater) updatePosts(blog database.Blog) {
	var receivedPosts []loaders.Post

	// Add posts
	posts, errors := u.loader.LoadPosts()
	for {
		select {
		case err, ok := <-errors:
			if ok {
				u.log("loader error: %s", err)
				u.stats.LoaderErrors++
			} else {
				errors = nil
			}
		case post, ok := <-posts:
			if ok {
				receivedPosts = append(receivedPosts, post)
				u.stats.PostsReceived++
				u.handlePost(blog, post)
			} else {
				posts = nil
			}
		}
		if errors == nil && posts == nil {
			break
		}
	}

	if u.stats.LoaderErrors == 0 {
		u.deleteOldPosts(blog, receivedPosts)
	}
}

// deleteOldPosts removes the posts which are in the database but are not on
// the blog.
func (u *blogUpdater) deleteOldPosts(blog database.Blog, receivedPosts []loaders.Post) {
	dbPosts, err := getPostsForDeleteOldPosts(blog)
	if err != nil {
		u.log("deleteOldPosts - couldn't retrieve the posts: %s", err)
	}
	u.stats.PostRemovalsStarted = true
	for _, post := range dbPosts {
		if !existsOnBlog(receivedPosts, post) {
			u.stats.PostRemovalsAttempted++
			if _, err := database.DB.Exec("DELETE FROM post WHERE id=$1", post.Post.ID); err == nil {
				u.stats.PostRemovalsSucceeded++
			} else {
				u.log("deleteOldPosts - couldn't remove a post: %s", err)
			}
		}
	}
}

func existsOnBlog(receivedPosts []loaders.Post, post postWithCategory) bool {
	for _, receivedPost := range receivedPosts {
		if receivedPost.Id == post.InternalID &&
			receivedPost.Category == post.Name {
			return true
		}
	}
	return false
}

// handlePost updates a single post received from the loader.
func (u *blogUpdater) handlePost(blog database.Blog, loadedPost loaders.Post) {
	tx, err := database.DB.Beginx()
	if err != nil {
		u.log("handlePost - couldn't begin transaction: %s", err)
		return
	}

	category, err := getCategory(tx, blog, loadedPost.Category)
	if err != nil {
		u.log("handlePost - category wasn't retrieved: %s", err)
		err := tx.Rollback()
		if err != nil {
			u.log("handlePost - couldn't rollback the transaction: %s", err)
		}
		return
	}

	post, err := getPost(tx, category, loadedPost.Id)
	if err != nil {
		u.log("handlePost - post wasn't retrieved: %s", err)
		err := tx.Rollback()
		if err != nil {
			u.log("handlePost - couldn't rollback the transaction: %s", err)
		}
		return
	}

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

	if altered {
		if _, err := tx.NamedExec("UPDATE post SET date=:date, title=:title, summary=:summary WHERE id=:id", post); err != nil {
			u.log("handlePost - post wasn't saved: %s", err)
			err := tx.Rollback()
			if err != nil {
				u.log("handlePost - couldn't rollback the transaction: %s", err)
			}
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		u.log("handlePost - couldn't commit the transaction: %s", err)
		return
	}

	if altered {
		u.stats.PostsUpdated++
	} else {
		u.stats.PostsUnaltered++
	}

	tags, err := getTags(post)
	if err != nil {
		u.log("tags weren't retrieved: %s", err)
		u.stats.TagAddErrors++
		return
	} else {
		u.addNewTags(post, loadedPost.Tags, tags)
		u.removeOldTags(post, loadedPost.Tags, tags)
	}
}

// addNewTags ensures that all tags are present in the database.
func (u *blogUpdater) addNewTags(post database.Post, loadedTags []string, postTags []database.Tag) {
	for _, loadedTagName := range loadedTags {
		loadedTag, err := getTag(loadedTagName)
		if err != nil {
			u.log("tag wasn't retrieved: %s", err)
			u.stats.TagAddErrors++
			continue
		}
		if tag := u.findTag(loadedTag, postTags); tag == nil {
			if _, err := database.DB.Exec("INSERT INTO post_to_tag (post_id, tag_id) VALUES ($1, $2)", post.ID, loadedTag.ID); err != nil {
				u.log("tag wasn't added: %s", err)
				u.stats.TagAddErrors++
				continue
			}
		}
	}
}

func (u *blogUpdater) findTag(loadedTag database.Tag, tags []database.Tag) *database.Tag {
	for _, tag := range tags {
		if tag.ID == loadedTag.ID {
			return &tag
		}
	}
	return nil
}

// removeOldTags ensures that all old tags are removed from the database.
func (u *blogUpdater) removeOldTags(post database.Post, loadedTags []string, postTags []database.Tag) {
	for _, tag := range postTags {
		if !u.findLoadedTag(tag.Name, loadedTags) {
			if _, err := database.DB.Exec("DELETE FROM post_to_tag WHERE id=$1", tag.ID); err != nil {
				u.log("tag wasn't removed: %s", err)
				u.stats.TagRemoveErrors++
				continue
			}
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

func (u *blogUpdater) log(format string, v ...interface{}) {
	log.Printf(u.loader.GetUrl()+": "+format, v...)
}
