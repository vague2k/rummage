package main

import (
	"fmt"
	"time"

	"github.com/vague2k/rummage/pkg/database"
)

func main() {
	r, err := database.Init("")
	if err != nil {
		panic(err)
	}
	defer r.DB.Close()

	_, err = r.AddItem("somecontent")
	if err != nil {
		panic(err)
	}

	item, _ := r.SelectItem("somecontent")
	fmt.Println(item)

	update := &database.RummageDBItem{
		Entry:        "updatedsomeothercontent",
		Score:        2.0,
		LastAccessed: time.Now().Unix(),
	}
	_, err = r.UpdateItem("somecontent", update)
	if err != nil {
		panic(err)
	}
}
