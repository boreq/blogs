package core

import (
	blgs "github.com/boreq/blogs/blogs"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type postResult struct {
	database.Post
	CategoryName   string
	BlogInternalID uint
}

func (p postResult) GetUrl() string {
	loader, ok := blgs.Blogs[p.BlogInternalID]
	if ok {
		return loader.GetPostUrl(p.Post.InternalID)
	}
	return ""
}

func (p postResult) GetTags() []database.Tag {
	var tags []database.Tag
	err := database.DB.Select(&tags,
		`SELECT tag.*
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id = tag.id
		JOIN post ON post.id = post_to_tag.post_id
		WHERE post.id=$1
		ORDER BY tag.name DESC`, p.Post.ID)
	if err != nil {
		panic(err)
	}
	return tags
}

type TagResult struct {
	database.Tag
	Count uint
}

func blog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the data
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		errors.NotFound(w, r)
		return
	}

	var blog database.Blog
	err = database.DB.Get(&blog, "SELECT * FROM blog WHERE id=$1", id)
	if err != nil {
		if err == database.ErrNoRows {
			errors.NotFound(w, r)
			return
		} else {
			errors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var subscription *database.Subscription
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user_id := ctx.User.GetUser().ID
		err = database.DB.Get(subscription,
			`SELECT * FROM
			subscription WHERE
			blog_id=$1 AND user_id=$2
			LIMIT 1`, id, user_id)
		if err != nil && err != database.ErrNoRows {
			errors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var categories []database.Category
	err = database.DB.Select(&categories,
		`SELECT category.*
		FROM category
		JOIN blog ON blog.id=category.blog_id
		WHERE blog.id=$1`, id)
	if err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var posts []postResult
	err = database.DB.Select(&posts,
		`SELECT post.*, category.name as category_name, blog.internal_id AS blog_internal_id
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		WHERE blog.id=$1
		ORDER BY post.date DESC`, id)
	if err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var tags []TagResult
	err = database.DB.Select(&tags,
		`SELECT tag.*, COUNT(tag.id) as count
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id = tag.id
		JOIN post ON post.id = post_to_tag.post_id
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		WHERE blog.id=$1
		GROUP BY tag.id
		ORDER BY count DESC`, id)
	if err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}

	// Render
	var data = templates.GetDefaultData(r)
	data["blog"] = blog
	data["subscription"] = subscription
	data["categories"] = categories
	data["posts"] = posts
	data["tags"] = tags
	if err := templates.RenderTemplateSafe(w, "core/blog.tmpl", data); err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}
}
