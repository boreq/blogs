// Package commands defines the commands used by the guinea library.
package commands

import (
	"errors"
	"github.com/boreq/blogs/config"
	"github.com/boreq/blogs/database"
)

func configInit(configFilename string) error {
	return config.Load(configFilename)
}

func coreInit(configFilename string) error {
	// Config
	if err := configInit(configFilename); err != nil {
		return err
	}

	// Database
	var dbType database.DatabaseType
	switch config.Config.DatabaseType {
	case "sqlite":
		dbType = database.SQLite3
		break
	case "postgresql":
		dbType = database.PostgreSQL
		break
	default:
		return errors.New("The legal values for DatabaseType config key are: \"sqlite\", \"postgresql\"")
	}
	if err := database.Init(dbType, config.Config.DatabaseURI); err != nil {
		return err
	}
	return nil
}
