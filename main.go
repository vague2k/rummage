package main

import (
	"github.com/vague2k/rummage/cmd"
	"github.com/vague2k/rummage/pkg/database"
)

func main() {
	db, err := database.Init("")
	if err != nil {
		panic(err)
	}

	root := cmd.NewRootCmd(db)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
