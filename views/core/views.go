package core

import (
	"github.com/boreq/blogs/database"
	bhttp "github.com/boreq/blogs/http"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/utils"
	verrors "github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var posts []result

	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user_id := ctx.User.GetUser().ID

		var numPosts uint
		if err := database.DB.Get(&numPosts,
			`SELECT COUNT(*) AS numPosts
			FROM post
			JOIN category ON category.id = post.category_id
			JOIN blog ON blog.id = category.blog_id
			JOIN subscription ON blog.id = subscription.blog_id
			WHERE subscription.user_id=$1`, user_id); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}

		p := utils.NewPagination(r, numPosts, 20)

		if err := database.DB.Select(&posts,
			`SELECT post.*, category.*, blog.*, subscription.*
			FROM post
			JOIN category ON category.id = post.category_id
			JOIN blog ON blog.id = category.blog_id
			JOIN subscription ON blog.id = subscription.blog_id
			WHERE subscription.user_id=$1
			ORDER BY post.date DESC
			LIMIT $2 OFFSET $3
			`, user_id, p.Limit, p.Offset); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var newPosts []postCategoryBlog
	if err := database.DB.Select(&newPosts,
		`SELECT post.*, category.*, blog.*
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		ORDER BY post.date DESC LIMIT 5`); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["posts"] = posts
	data["new_posts"] = newPosts
	if err := templates.RenderTemplateSafe(w, "core/index.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func posts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var numPosts uint
	if err := database.DB.Get(&numPosts,
		`SELECT COUNT(*) AS numPosts
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id`); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
	p := utils.NewPagination(r, numPosts, 20)

	var posts []postCategoryBlog
	if err := database.DB.Select(&posts,
		`SELECT post.*, category.*, blog.*
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		ORDER BY post.date DESC
		LIMIT $1 OFFSET $2`, p.Limit, p.Offset); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["posts"] = posts
	if err := templates.RenderTemplateSafe(w, "core/posts.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func blogs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the data
	user_id := -1
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user_id = int(ctx.User.GetUser().ID)
	}

	var blogs = make([]blogResult, 0)
	err := database.DB.Select(&blogs, `
		SELECT blog.*, subscription.id as subscription_id, MAX(post.date) AS updated
		FROM blog
		JOIN category ON category.blog_id = blog.id
		JOIN post ON post.category_id = category.id
		LEFT JOIN subscription ON subscription.blog_id = blog.id AND subscription.user_id=$1
		GROUP BY blog.id
		ORDER BY blog.title`, user_id)
	if err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["blogs"] = blogs
	if err := templates.RenderTemplateSafe(w, "core/blogs.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func blog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the data
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		verrors.NotFound(w, r)
		return
	}

	var blog database.Blog
	err = database.DB.Get(&blog, "SELECT * FROM blog WHERE id=$1", id)
	if err != nil {
		if err == database.ErrNoRows {
			verrors.NotFound(w, r)
			return
		} else {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	subscription := &database.Subscription{}
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user_id := ctx.User.GetUser().ID
		err = database.DB.Get(subscription,
			`SELECT * FROM
			subscription WHERE
			blog_id=$1 AND user_id=$2
			LIMIT 1`, id, user_id)
		if err != nil {
			if err == database.ErrNoRows {
				subscription = nil
			} else {
				verrors.InternalServerErrorWithStack(w, r, err)
				return
			}
		}
	}

	var categories []database.Category
	err = database.DB.Select(&categories,
		`SELECT category.*
		FROM category
		JOIN blog ON blog.id=category.blog_id
		WHERE blog.id=$1`, id)
	if err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
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
		verrors.InternalServerErrorWithStack(w, r, err)
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
		verrors.InternalServerErrorWithStack(w, r, err)
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
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func subscribe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blog_id, err := strconv.ParseUint(r.FormValue("blog_id"), 10, 32)
	if err != nil {
		verrors.BadRequest(w, r)
		return
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		bhttp.RedirectOrNext(w, r, "/")
		return
	}
	user_id := ctx.User.GetUser().ID

	if _, err := database.DB.Exec(`
		INSERT INTO subscription(blog_id, user_id) 
		SELECT $1, $2
		WHERE NOT EXISTS(
			SELECT 1
			FROM subscription
			WHERE blog_id=$1 AND user_id=$2)`,
		blog_id, user_id); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	bhttp.RedirectOrNext(w, r, "/")
}

func unsubscribe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blog_id, err := strconv.ParseUint(r.FormValue("blog_id"), 10, 32)
	if err != nil {
		verrors.BadRequest(w, r)
		return
	}

	ctx := context.Get(r)
	if !ctx.User.IsAuthenticated() {
		bhttp.RedirectOrNext(w, r, "/")
		return
	}
	user_id := ctx.User.GetUser().ID

	if _, err := database.DB.Exec(`
		DELETE FROM subscription
		WHERE blog_id=$1 AND user_id=$2`,
		blog_id, user_id); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	bhttp.RedirectOrNext(w, r, "/")
}
