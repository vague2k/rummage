package main

import (
	"fmt"
	"os"

	"github.com/vague2k/rummage/pkg/db"
)

func main() {
	db, err := db.Access()
	if err != nil {
		panic(err)
	}
	err = db.AddItem("someothercontent")
	if err != nil {
		panic(err)
	}

	b, err := os.ReadFile(db.FilePath)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}
