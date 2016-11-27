package commands

import (
	"github.com/boreq/blogs/blogs/updater"
	"github.com/boreq/blogs/config"
	"github.com/boreq/guinea"
	"strconv"
)

var testLoaderCmd = guinea.Command{
	Run: runTestLoader,
	Arguments: []guinea.Argument{
		{"config", false, "Config file"},
		{"id", false, "ID of the loader"},
	},
	ShortDescription: "tests a loader without altering the database",
}

func runTestLoader(c guinea.Context) error {
	if err := configInit(c.Arguments[0]); err != nil {
		return err
	}
	config.Config.Debug = true
	id, err := strconv.Atoi(c.Arguments[1])
	if err != nil {
		return err
	}
	return updater.TestLoader(uint(id))
}
