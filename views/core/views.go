package core

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/utils"
	verrors "github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

const postsPerPage = 20
const blogsPerPage = 20
const updatesPerPage = 30

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s := utils.NewSort(r, []utils.SortParam{
		{Key: "date", Label: "Date", Query: "post.date", Reversed: true},
		{Key: "stars", Label: "Stars", Query: "post.stars", Reversed: true},
		{Key: "title", Label: "Title", Query: "post.title"},
	})
	preserveParams := make(map[string]string)
	preserveParams["sort"] = s.CurrentKey
	var posts []postsResult

	ctx := context.Get(r)
	var pagination utils.Pagination
	if ctx.User.IsAuthenticated() {
		userId := ctx.User.GetUser().ID
		var numPosts uint
		if err := database.DB.Get(&numPosts,
			`SELECT COUNT(*) AS numPosts
			FROM post
			JOIN category ON category.id = post.category_id
			JOIN blog ON blog.id = category.blog_id
			JOIN subscription ON blog.id = subscription.blog_id
			WHERE subscription.user_id=$1`, userId); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
		pagination = utils.NewPagination(r, numPosts, postsPerPage, preserveParams)
		if err := database.DB.Select(&posts,
			`SELECT post.*, category.*, blog.*, star.id AS starred
			FROM post
			JOIN category ON category.id = post.category_id
			JOIN blog ON blog.id = category.blog_id
			JOIN subscription ON blog.id = subscription.blog_id
			LEFT JOIN star ON star.post_id=post.id AND star.user_id=$1
			WHERE subscription.user_id=$1
			ORDER BY `+s.Query+`
			LIMIT $2 OFFSET $3
			`, userId, pagination.Limit, pagination.Offset); err != nil {
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
	data["pagination"] = pagination
	data["sort"] = s
	if err := templates.RenderTemplateSafe(w, "core/index.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func posts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// tag
	tag := ""
	tagParams, ok := r.URL.Query()["tag"]
	if ok {
		tag = tagParams[0]
	}

	// userId
	var userId uint
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = ctx.User.GetUser().ID
	}

	// Sorting
	s := utils.NewSort(r, []utils.SortParam{
		{Key: "date", Label: "Date", Query: "post.date", Reversed: true},
		{Key: "stars", Label: "Stars", Query: "post.stars", Reversed: true},
		{Key: "title", Label: "Title", Query: "post.title"},
	})
	preserveParams := make(map[string]string)
	preserveParams["sort"] = s.CurrentKey
	if tag != "" {
		preserveParams["tag"] = tag
	}

	// Prepare queries
	var numPostsQuery string
	var postsQuery string
	if tag != "" {
		tagJoin := `
		JOIN post_to_tag ON post_to_tag.post_id=post.id
		JOIN tag ON post_to_tag.tag_id=tag.id
		`
		tagWhere := "WHERE tag.name=$2"
		postsQuery = `
		SELECT post.*, category.*, blog.*, star.id AS starred
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		` + tagJoin + `
		LEFT JOIN star ON star.post_id=post.id AND star.user_id=$1
		` + tagWhere + `
		GROUP BY post.id
		ORDER BY ` + s.Query + `
		LIMIT $3 OFFSET $4`
		numPostsQuery = `SELECT COUNT(*) AS numPosts FROM post
			JOIN category ON category.id = post.category_id
			JOIN blog ON blog.id = category.blog_id
			` + tagJoin + " " + tagWhere
	} else {
		postsQuery = `SELECT post.*, category.*, blog.*, star.id AS starred
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		LEFT JOIN star ON star.post_id=post.id AND star.user_id=$1
		GROUP BY post.id
		ORDER BY ` + s.Query + `
		LIMIT $2 OFFSET $3`
		numPostsQuery = "SELECT COUNT(*) AS numPosts FROM post"
	}

	// Execute
	var numPosts uint
	var posts []postsResult

	if tag == "" {
		if err := database.DB.Get(&numPosts, numPostsQuery); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	} else {
		if err := database.DB.Get(&numPosts, numPostsQuery, tag); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	p := utils.NewPagination(r, numPosts, postsPerPage, preserveParams)

	if tag == "" {
		if err := database.DB.Select(&posts, postsQuery, userId, p.Limit, p.Offset); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	} else {
		if err := database.DB.Select(&posts, postsQuery, userId, tag, p.Limit, p.Offset); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var data = templates.GetDefaultData(r)
	data["posts"] = posts
	data["pagination"] = p
	data["sort"] = s
	data["tag"] = tag
	if err := templates.RenderTemplateSafe(w, "core/posts.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func blogs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Sorting
	s := utils.NewSort(r, []utils.SortParam{
		{Key: "title", Label: "Title", Query: "blog.title"},
		{Key: "subscribers", Label: "Subscribers", Query: "blog.subscriptions", Reversed: true},
		{Key: "last_post", Label: "Last post", Query: "updated", Reversed: true},
	})
	preserveParams := map[string]string{
		"sort": s.CurrentKey,
	}

	// User id for displaying subscription buttons
	userId := -1
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		userId = int(ctx.User.GetUser().ID)
	}

	// Count
	var numBlogs uint
	if err := database.DB.Get(&numBlogs, `SELECT COUNT(*) AS numBlogs FROM blog`); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
	p := utils.NewPagination(r, numBlogs, blogsPerPage, preserveParams)

	// Query
	var blogs []blogResult
	if numBlogs > 0 {
		if err := database.DB.Select(&blogs,
			`SELECT blog.*, MAX(post.date) AS updated, subscription.id AS subscribed
			FROM blog
			JOIN category ON category.blog_id=blog.id
			JOIN post ON post.category_id=category.id
			LEFT JOIN subscription ON subscription.blog_id=blog.id AND subscription.user_id=$1
			GROUP BY blog.id
			ORDER BY `+s.Query+`
			LIMIT $2 OFFSET $3`, userId, p.Limit, p.Offset); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	// Render templates
	var data = templates.GetDefaultData(r)
	data["blogs"] = blogs
	data["sort"] = s
	data["pagination"] = p
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
	var user_id uint
	if ctx.User.IsAuthenticated() {
		user_id = ctx.User.GetUser().ID
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

	var posts []postsResult
	err = database.DB.Select(&posts,
		`SELECT post.*, category.*, blog.*, star.id AS starred
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		LEFT JOIN star ON star.post_id=post.id AND star.user_id=$1
		WHERE blog.id=$2
		GROUP BY post.id
		ORDER BY post.date DESC`, user_id, id)
	if err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var tags []tagResult
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

func profile_stars(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		verrors.BadRequest(w, r)
		return
	}

	var profile database.User
	err = database.DB.Get(&profile, "SELECT * FROM user WHERE id=$1", userId)
	if err != nil {
		if err == database.ErrNoRows {
			verrors.NotFound(w, r)
			return
		} else {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var numPosts uint
	if err := database.DB.Get(&numPosts,
		`SELECT COUNT(*) AS numPosts
		FROM post
		JOIN star ON star.post_id = post.id
		JOIN user ON user.id = star.user_id
		WHERE user.id=$1`, userId); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
	pagination := utils.NewPagination(r, numPosts, postsPerPage, nil)
	var posts []postsResult
	if err := database.DB.Select(&posts,
		`SELECT post.*, category.*, blog.*, star.id AS starred
		FROM post
		JOIN category ON category.id = post.category_id
		JOIN blog ON blog.id = category.blog_id
		JOIN star ON star.post_id = post.id
		JOIN user ON user.id = star.user_id
		WHERE user.id=$1
		LIMIT $2 OFFSET $3
			`, userId, pagination.Limit, pagination.Offset); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["profile"] = profile
	data["posts"] = posts
	data["pagination"] = pagination
	if err := templates.RenderTemplateSafe(w, "core/profile_stars.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func profile_subscriptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		verrors.BadRequest(w, r)
		return
	}

	var profile database.User
	err = database.DB.Get(&profile, "SELECT * FROM user WHERE id=$1", userId)
	if err != nil {
		if err == database.ErrNoRows {
			verrors.NotFound(w, r)
			return
		} else {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var numBlogs uint
	if err := database.DB.Get(&numBlogs,
		`SELECT COUNT(*) AS numBlogs
		FROM blog
		JOIN subscription ON subscription.blog_id = blog.id
		JOIN user ON user.id = subscription.user_id
		WHERE user.id=$1`, userId); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
	pagination := utils.NewPagination(r, numBlogs, postsPerPage, nil)
	var blogs []blogResult
	if numBlogs > 0 {
		if err := database.DB.Select(&blogs,
			`SELECT blog.*, MAX(post.date) AS updated, subscription.id AS subscribed
		FROM blog
		JOIN category ON category.blog_id=blog.id
		JOIN post ON post.category_id=category.id
		JOIN subscription ON subscription.blog_id = blog.id
		JOIN user ON user.id = subscription.user_id
		WHERE user.id=$1
		LIMIT $2 OFFSET $3
			`, userId, pagination.Limit, pagination.Offset); err != nil {
			verrors.InternalServerErrorWithStack(w, r, err)
			return
		}
	}

	var data = templates.GetDefaultData(r)
	data["profile"] = profile
	data["blogs"] = blogs
	data["pagination"] = pagination
	if err := templates.RenderTemplateSafe(w, "core/profile_subscriptions.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func tags(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s := utils.NewSort(r, []utils.SortParam{
		{Key: "name", Label: "Name", Query: "tag.name"},
		{Key: "posts", Label: "Posts", Query: "count", Reversed: true},
	})
	preserveParams := make(map[string]string)
	preserveParams["sort"] = s.CurrentKey
	var numTags uint
	if err := database.DB.Get(&numTags, `SELECT COUNT(*) AS numTags FROM tag`); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
	p := utils.NewPagination(r, numTags, 28, preserveParams)
	var tags []tagResult
	if err := database.DB.Select(&tags,
		`SELECT tag.*, COUNT(*) AS count
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id=tag.id
		GROUP BY tag.id
		ORDER BY `+s.Query+`
		LIMIT $1 OFFSET $2`, p.Limit, p.Offset); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["tags"] = tags
	data["pagination"] = p
	data["sort"] = s
	if err := templates.RenderTemplateSafe(w, "core/tags.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func updates(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s := utils.NewSort(r, []utils.SortParam{
		{Key: "date", Label: "Date", Query: "\"update\".started", Reversed: true},
	})
	f := utils.NewFilter(r, []utils.FilterParam{
		{Key: "all", Label: "All", Query: ""},
		{Key: "success", Label: "Succeeded", Query: "\"update\".succeeded=1"},
		{Key: "failure", Label: "Failed", Query: "\"update\".succeeded=0"},
	})
	preserveParams := make(map[string]string)
	preserveParams["sort"] = s.CurrentKey
	preserveParams["filter"] = f.CurrentKey

	where := ""
	if f.Query != "" {
		where = "WHERE " + f.Query
	}

	// Count and paginate
	var numUpdates uint
	if err := database.DB.Get(&numUpdates, `SELECT COUNT(*) AS numUpdates FROM "update"`+where); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
	p := utils.NewPagination(r, numUpdates, updatesPerPage, preserveParams)

	// Get rows
	var updates []updateResult
	if err := database.DB.Select(&updates,
		`SELECT "update".*, blog.*
		FROM "update"
		JOIN blog ON blog.id="update".blog_id
		`+where+`
		ORDER BY `+s.Query+`
		LIMIT $1 OFFSET $2`, p.Limit, p.Offset); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["updates"] = updates
	data["pagination"] = p
	data["sort"] = s
	data["filter"] = f
	if err := templates.RenderTemplateSafe(w, "core/updates.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}
