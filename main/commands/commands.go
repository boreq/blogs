// Package commands defines the commands used by the guinea library.
package commands

import (
	"github.com/boreq/blogs/config"
	"github.com/boreq/blogs/database"
)

func configInit(configFilename string) error {
	return config.Load(configFilename)
}

func coreInit(configFilename string) error {
	if err := configInit(configFilename); err != nil {
		return err
	}
	if err := database.Init(database.SQLite3, config.Config.DatabaseURI); err != nil {
		return err
	}
	return nil
}
