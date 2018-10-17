package commands

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/guinea"
)

var dropDbCmd = guinea.Command{
	Run: runDropDb,
	Arguments: []guinea.Argument{
		{"config", false, "Config file"},
	},
	ShortDescription: "drops database tables",
}

func runDropDb(c guinea.Context) error {
	if err := coreInit(c.Arguments[0]); err != nil {
		return err
	}
	return database.DropTables()
}
