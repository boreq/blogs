package commands

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/guinea"
)

var migrateDbCmd = guinea.Command{
	Run: runMigrateDb,
	Arguments: []guinea.Argument{
		{"config", false, "Config file"},
	},
	ShortDescription: "creates missing tables and columns in the database",
}

func runMigrateDb(c guinea.Context) error {
	if err := coreInit(c.Arguments[0]); err != nil {
		return err
	}
	return database.MigrateTables()
}
