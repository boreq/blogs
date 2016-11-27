package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Config *ConfigStruct = createDefaultConfig()

type ConfigStruct struct {
	Debug       bool
	DatabaseURI string
}

func createDefaultConfig() *ConfigStruct {
	conf := &ConfigStruct{
		Debug:       false,
		DatabaseURI: "/tmp/database.sqlite3",
	}
	return conf
}

// Load loads the config from the specified json file. If the certain keys
// are not present in the loaded config file, the default values are used.
func Load(filename string) error {
	content, e := ioutil.ReadFile(filename)
	if os.IsNotExist(e) {
		return nil
	}
	return json.Unmarshal(content, Config)
}
