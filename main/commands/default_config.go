package commands

import (
	"encoding/json"
	"fmt"
	"github.com/boreq/blogs/config"
	"github.com/boreq/guinea"
)

var defaultConfigCmd = guinea.Command{
	Run:              runDefaultConfig,
	ShortDescription: "prints the default configuration to stdout",
	Description: `
This command prints out the default config in the json format to stdout. The
available config keys are:

Debug
	Specifies if the program should run in the debug mode. The program
	running in the debug mode prints more log messages.
	Allowed values: true or false.

DatabaseURI
	Database-specific connection information.
	Allowed values:
		For sqlite:
			A string, path to a database file.
			Example: "/tmp/database.sqlite3"
		For postgresql:
			See https://godoc.org/github.com/lib/pq.
			Example: "postgres://user:password@localhost/database"

DatabaseType
	The type of the database to use.
	Allowed values: "sqlite" or "postgresql".

TemplatesDirectory
	A path of the templates directory. The path may be relative to the cwd
	or absolute.
	Allowed values: a string.

StaticDirectory
	A path of the directory containing the static files. The path may be
	relative to the cwd or absolute.
	Allowed values: a string.


	`,
}

func runDefaultConfig(c guinea.Context) error {
	defaultConfig := config.Default()
	j, err := json.MarshalIndent(defaultConfig, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(j))
	return nil
}
