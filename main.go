package main

import (
	"fmt"

	"github.com/vague2k/rummage/pkg/db"
)

func main() {
	db, err := db.Access("")
	if err != nil {
		panic(err)
	}
	item, err := db.AddItem("someothercontent")
	if err != nil {
		panic(err)
	}

	item = item.RecalculateScore()

	item, err = db.UpdateItem(item.Entry, item)
	if err != nil {
		panic(err)
	}
	// b, err := os.ReadFile(db.FilePath)
	// fmt.Print(string(b))

	fmt.Println("UPDATED ENTRY: ", item.Entry)
	fmt.Println("UPDATED SCORE: ", item.Score)
	fmt.Println("UPDATED EPOCH:", item.LastAccessed)
}
