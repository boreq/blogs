package core

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/templates"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type TagResult struct {
	database.Tag
	Count uint
}

func blog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	// Get the data
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		return err
	}

	var blog = &database.Blog{}
	database.DB.Where(&database.Blog{ID: uint(id)}).First(&blog)

	var categories = make([]database.Category, 0)
	database.DB.Model(&blog).Related(&categories)

	var posts = make([]database.Post, 0)
	database.DB.
		Order("posts.date desc").
		Joins("JOIN categories ON categories.id = posts.category_id JOIN blogs ON blogs.id = categories.blog_id").
		Where("blogs.id = ?", id).
		Preload("Tags").
		Find(&posts)

	var tags = make([]TagResult, 0)
	database.DB.
		Table("tags").
		Order("count desc").
		Select("tags.name, COUNT(tags.name) as count").
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id JOIN posts ON posts.id = post_tags.post_id JOIN categories ON categories.id = posts.category_id JOIN blogs ON blogs.id = categories.blog_id").
		Where("blogs.id = ?", id).
		Group("tags.name").
		Scan(&tags)

	// Render
	var data = make(map[string]interface{})
	data["blog"] = blog
	data["categories"] = categories
	data["posts"] = posts
	data["tags"] = tags
	return templates.RenderTemplate(w, "core/blog.tmpl", data)
}
