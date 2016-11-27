package commands

import (
	"github.com/boreq/blogs/config"
	"github.com/boreq/blogs/database"
)

func coreInit(configFilename string) error {
	if err := config.Load(configFilename); err != nil {
		return err
	}
	if err := database.Init(database.SQLite3, config.Config.DatabaseURI); err != nil {
		return err
	}
	return nil
}
