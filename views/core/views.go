package core

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/templates"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

type BlogWithDate struct {
	database.Blog
	Updated time.Time
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	blogs := make([]BlogWithDate, 0)
	//database.DB.Find(&blogs)

	// doesn't work
	database.DB.Table("blogs").Order("blogs.title").Select("blogs.id, blogs.internal_id, blogs.title, MAX(posts.date) AS updated").Joins("JOIN categories ON categories.blog_id = blogs.id JOIN posts ON posts.category_id = categories.id").Group("blogs.id").Scan(&blogs)

	var data map[string]interface{} = make(map[string]interface{})
	data["blogs"] = blogs
	return templates.RenderTemplate(w, "core/index.tmpl", data)
}

type TagsWithCount struct {
	database.Tag
	Count uint
}

func blog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		return err
	}

	blog := &database.Blog{}
	database.DB.Where(&database.Blog{ID: uint(id)}).First(&blog)

	categories := make([]database.Category, 0)
	database.DB.Model(&blog).Related(&categories)

	posts := make([]database.Post, 0)
	database.DB.Order("posts.date desc").Joins("JOIN categories ON categories.id = posts.category_id JOIN blogs ON blogs.id = categories.blog_id").Where("blogs.id = ?", id).Preload("Tags").Find(&posts)

	tags := make([]TagsWithCount, 0)
	database.DB.Table("tags").Order("count desc").Select("tags.name, COUNT(tags.name) as count").Order("posts.date desc").Joins("JOIN post_tags ON post_tags.tag_id = tags.id JOIN posts ON posts.id = post_tags.post_id JOIN categories ON categories.id = posts.category_id JOIN blogs ON blogs.id = categories.blog_id").Where("blogs.id = ?", id).Group("tags.name").Scan(&tags)
	//db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)

	var data map[string]interface{} = make(map[string]interface{})
	data["blog"] = blog
	data["categories"] = categories
	data["posts"] = posts
	data["tags"] = tags
	return templates.RenderTemplate(w, "core/blog.tmpl", data)
}
