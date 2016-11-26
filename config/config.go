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

func Load(filename string) error {
	content, e := ioutil.ReadFile(filename)
	if os.IsNotExist(e) {
		return nil
	}
	return json.Unmarshal(content, Config)
}
