package main

import (
	"net/http"

	"shop/app"
	"shop/db"
)

func main() {
	application := &app.App{
		DB:  db.Connect(),
		Key: []byte("eifbwrrpbep"),
	}

	http.HandleFunc("/login", application.LoginHandler)
	http.HandleFunc("/product/add", application.AddProduct)
	http.HandleFunc("/products", application.GetProducts)
	http.HandleFunc("/walet", application.AddWalet)
	http.HandleFunc("/product/select", application.SelectProduct)
	http.HandleFunc("/cart", application.ShowCart)

	http.ListenAndServe(":8080", nil)
}
