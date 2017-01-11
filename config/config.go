// Package config holds the global config.
package config

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigStruct struct {
	Debug              bool
	DatabaseURI        string
	DatabaseType       string
	TemplatesDirectory string
	StaticDirectory    string
}

// Config points to the current config struct used by the other parts of the
// program.
var Config *ConfigStruct = Default()

// Default returns the default config.
func Default() *ConfigStruct {
	conf := &ConfigStruct{
		Debug:              false,
		DatabaseURI:        "/tmp/database.sqlite3",
		DatabaseType:       "sqlite",
		TemplatesDirectory: "_templates",
		StaticDirectory:    "_static",
	}
	return conf
}

// Load loads the config from the specified json file. If certain keys are not
// present in the loaded config file the current values of the config struct
// are preserved.
func Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, Config)
}
