package db

import (
	"context"
	"fmt"

	"shop/ent"

	_ "github.com/mattn/go-sqlite3"
)

var Client *ent.Client
var Err error

func Connect() {
	Client, Err = ent.Open("sqlite3", "file:shop.db?_fk=1")
	if Err != nil {
		fmt.Println(Err)
	}
	Client.Schema.Create(context.Background())
}
func Close() {
	if Client != nil {
		Client.Close()
	}
}
