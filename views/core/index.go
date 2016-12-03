package core

import (
	blgs "github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type result struct {
	database.Post
	database.Category
	database.Blog
	database.Subscription
}

func (p result) GetUrl() string {
	loader, ok := blgs.Blogs[p.Blog.InternalID]
	if ok {
		return loader.GetPostUrl(p.Post.InternalID)
	}
	return ""
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var posts []result

	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user_id := ctx.User.GetUser().ID
		err := database.DB.Select(&posts,
			`SELECT post.*, category.*, blog.*, subscription.*
			FROM post
			JOIN category ON category.id = post.category_id
			JOIN blog ON blog.id = category.blog_id
			JOIN subscription ON blog.id = subscription.blog_id
			WHERE subscription.user_id=$1
			ORDER BY post.date DESC`, user_id)
		if err != nil {
			errors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var data = templates.GetDefaultData(r)
	data["posts"] = posts
	if err := templates.RenderTemplateSafe(w, "core/index.tmpl", data); err != nil {
		errors.InternalServerError(w, r)
		return
	}
}
