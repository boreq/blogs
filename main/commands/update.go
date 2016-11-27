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
	Options: []guinea.Option{
		guinea.Option{
			Name:        "id",
			Type:        guinea.Int,
			Description: "ID of the blog to update, updates all if not present",
			Default:     -1,
		},
	},
	ShortDescription: "loads the data",
}

func runUpdate(c guinea.Context) error {
	if err := coreInit(c.Arguments[0]); err != nil {
		return err
	}
	id := c.Options["id"].Int()
	if id != -1 {
		return updater.UpdateSpecific(uint(id))
	} else {
		return updater.Update()
	}
}
