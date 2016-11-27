package commands

import (
	"github.com/boreq/blogs/blogs/updater"
	"github.com/boreq/guinea"
)

var updateCmd = guinea.Command{
	Run: runUpdate,
	Arguments: []guinea.Argument{
		{"config", false, "Config file"},
	},
	ShortDescription: "loads the data",
}

func runUpdate(c guinea.Context) error {
	if err := coreInit(c.Arguments[0]); err != nil {
		return err
	}
	return updater.Update()
}
