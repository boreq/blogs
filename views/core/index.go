package core

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type BlogResult struct {
	database.Blog
	Updated []uint8
}

func (br BlogResult) GetUpdated() string {
	t, err := time.Parse("2006-01-02 15:04:05-07:00", string(br.Updated))
	if err != nil {
		return ""
	}
	return utils.ISO8601(t)
}

type BlogData struct {
	database.Blog
	Updated string
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	// Get the data
	var blogs = make([]BlogResult, 0)
	database.DB.
		Table("blogs").
		Order("blogs.title").
		Select("blogs.id, blogs.internal_id, blogs.title, MAX(posts.date) AS updated").Joins("JOIN categories ON categories.blog_id = blogs.id JOIN posts ON posts.category_id = categories.id").
		Group("blogs.id").
		Scan(&blogs)

	// Render
	var blogsData = make([]BlogData, 0, len(blogs))
	for _, blog := range blogs {
		blogData := BlogData{blog.Blog, blog.GetUpdated()}
		blogsData = append(blogsData, blogData)
	}
	var data = make(map[string]interface{})
	data["blogs"] = blogsData
	return templates.RenderTemplate(w, "core/index.tmpl", data)
}
