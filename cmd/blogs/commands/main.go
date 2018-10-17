package commands

import (
	"fmt"
	"github.com/boreq/guinea"
)

var buildCommit string
var buildDate string

var MainCmd = guinea.Command{
	Options: []guinea.Option{
		guinea.Option{
			Name:        "version",
			Type:        guinea.Bool,
			Description: "Display version",
		},
	},
	Run: runMain,
	Subcommands: map[string]*guinea.Command{
		"serve":          &serveCmd,
		"update":         &updateCmd,
		"create_tables":  &createDbCmd,
		"drop_tables":    &dropDbCmd,
		"test_loader":    &testLoaderCmd,
		"list_loaders":   &listLoadersCmd,
		"default_config": &defaultConfigCmd,
	},
	ShortDescription: "a blog aggregation platform",
	Description:      "A web-based blog aggregation platform.",
}

func runMain(c guinea.Context) error {
	if c.Options["version"].Bool() {
		fmt.Printf("BuildCommit %s\n", buildCommit)
		fmt.Printf("BuildDate %s\n", buildDate)
		return nil
	}
	return guinea.ErrInvalidParms
}
