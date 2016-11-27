package commands

import (
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/views"
	"github.com/boreq/guinea"
	"github.com/julienschmidt/httprouter"
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
	if err := templates.Load(); err != nil {
		return err
	}
	return nil
}

func runServe(c guinea.Context) error {
	if err := initServe(c.Arguments[0]); err != nil {
		return err
	}

	router := httprouter.New()
	views.Register(router)
	return http.ListenAndServe(":8080", context.ClearHandler(router))
}
