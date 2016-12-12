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
