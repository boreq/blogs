package commands

import (
	"github.com/boreq/blogs/config"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/http/handler"
	blogService "github.com/boreq/blogs/service/blog"
	blogsService "github.com/boreq/blogs/service/blogs"
	"github.com/boreq/blogs/views/auth"
	"github.com/boreq/blogs/views/blog"
	"github.com/boreq/blogs/views/blogs"
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

	blogs := blogs.New("/blogs", blogsService)
	blog := blog.New("/blog", blogService)
	auth := auth.New("/auth")

	registerers := []handler.Registerer{
		blogs,
		blog,
		auth,
	}

	return http.ListenAndServe(config.Config.ServeAddress, handler.New(registerers))
}
