package commands

import (
	"github.com/boreq/blogs/config"
	"github.com/boreq/blogs/database"
	bhttp "github.com/boreq/blogs/http"
	authService "github.com/boreq/blogs/service/auth"
	blogService "github.com/boreq/blogs/service/blog"
	blogsService "github.com/boreq/blogs/service/blogs"
	contextService "github.com/boreq/blogs/service/context"
	postsService "github.com/boreq/blogs/service/posts"
	tagService "github.com/boreq/blogs/service/tag"
	"github.com/boreq/blogs/views/auth"
	"github.com/boreq/blogs/views/blog"
	"github.com/boreq/blogs/views/blogs"
	"github.com/boreq/blogs/views/post"
	"github.com/boreq/blogs/views/posts"
	"github.com/boreq/guinea"
	"net/http"
)

var serveCmd = guinea.Command{
	Run: runServe,
	Arguments: []guinea.Argument{
		{"config", false, "Config file"},
	},
	ShortDescription: "runs a server",
}

func initServe(configFilename string) error {
	if err := coreInit(configFilename); err != nil {
		return err
	}
	return nil
}

func runServe(c guinea.Context) error {
	if err := initServe(c.Arguments[0]); err != nil {
		return err
	}

	blogsService := blogsService.New(database.DB)
	blogService := blogService.New(database.DB)
	postsService := postsService.New(database.DB)
	tagService := tagService.New(database.DB)
	authService := authService.New(database.DB)
	contextService := contextService.New(authService)

	auth := auth.New("/auth", authService, contextService)
	blogs := blogs.New("/blogs", blogsService, contextService)
	blog := blog.New("/blog", blogService, postsService, tagService, contextService)
	posts := posts.New("/posts", postsService, tagService, contextService)
	post := post.New("/post", postsService, contextService)

	registerers := []bhttp.Registerer{
		auth,
		blogs,
		blog,
		posts,
		post,
	}

	return http.ListenAndServe(config.Config.ServeAddress, bhttp.New(registerers, contextService))
}
