package commands

import (
	"fmt"
	"github.com/boreq/blogs/blogs"
	"github.com/boreq/guinea"
)

var listLoadersCmd = guinea.Command{
	Run:              runListLoaders,
	ShortDescription: "lists all loaders and their IDs",
}

func runListLoaders(c guinea.Context) error {
	for internalID, loader := range blogs.Blogs {
		fmt.Println(internalID, loader.GetUrl())
	}
	return nil
}
