package commands

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/guinea"
)

var createDbCmd = guinea.Command{
	Run: runCreateDb,
	Arguments: []guinea.Argument{
		{"config", false, "Config file"},
	},
	ShortDescription: "creates database tables",
}

func runCreateDb(c guinea.Context) error {
	if err := coreInit(c.Arguments[0]); err != nil {
		return err
	}
	return database.CreateTables()
}
