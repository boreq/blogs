package commands

import "github.com/boreq/guinea"

var MainCmd = guinea.Command{
	Options: []guinea.Option{
		guinea.Option{
			Name:        "version",
			Type:        guinea.Bool,
			Description: "Display version",
		},
	},
	Run: func(c guinea.Context) error {
		if c.Options["version"].Bool() {
			return nil
		}
		return guinea.ErrInvalidParms
	},
	Subcommands: map[string]*guinea.Command{
		"serve":            &serveCmd,
		"update":           &updateCmd,
		"drop_database":    &dropDbCmd,
		"migrate_database": &migrateDbCmd,
	},
	ShortDescription: "a blog aggregation platform",
	Description: `Main command decription.
Second line.`,
}
