package commands

import (
	"github.com/boreq/blogs/config"
	"github.com/boreq/blogs/http/handler"
	"github.com/boreq/blogs/templates"
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
	if err := templates.Load(config.Config.TemplatesDirectory); err != nil {
		return err
	}
	return nil
}

func runServe(c guinea.Context) error {
	if err := initServe(c.Arguments[0]); err != nil {
		return err
	}
	return http.ListenAndServe(":8080", handler.Get())
}
