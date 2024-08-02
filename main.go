package main

import (
	"github.com/vague2k/rummage/cmd"
)

func main() {
	// db, err := database.Init("")
	// if err != nil {
	//     panic(err)
	// }

	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
