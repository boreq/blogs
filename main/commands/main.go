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
		"serve":         &serveCmd,
		"update":        &updateCmd,
		"create_tables": &createDbCmd,
		"drop_tables":   &dropDbCmd,
		"test_loader":   &testLoaderCmd,
		"list_loaders":  &listLoadersCmd,
	},
	ShortDescription: "a blog aggregation platform",
	Description: `Main command decription.
Second line.`,
}
