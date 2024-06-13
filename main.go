package main

import (
	"fmt"
	"time"

	"github.com/vague2k/rummage/pkg/db"
)

func main() {
	r, err := db.InitRummageDB("")
	if err != nil {
		panic(err)
	}
	defer r.DB.Close()

	r.AddItem("somecontent")

	item, _ := r.SelectItem("somecontent")
	fmt.Println(item)

	update := &db.RummageDBItem{
		Entry:        "updatedsomeothercontent",
		Score:        2.0,
		LastAccessed: time.Now().Unix(),
	}
	r.UpdateItem("somecontent", update)
}
