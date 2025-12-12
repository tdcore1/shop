package db

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"shop/ent"
)

func Connect() *ent.Client {
	client, err := ent.Open("sqlite3", "file:shop.db?_fk=1")
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	return client
}
