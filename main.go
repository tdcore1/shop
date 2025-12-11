package main

import (
	"fmt"
	"net/http"

	"shop/db"
	"shop/handlers"
)

func main() {
	db.Connect()
	defer db.Close()
	http.HandleFunc("/product", handlers.HandlerProduct)
	http.HandleFunc("/product/add", handlers.HandlerAdd)
	http.HandleFunc("/user", handlers.HandlerUser)
	http.HandleFunc("/walet", handlers.HandlerWalet)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/product/select", handlers.CartHandler)
	http.HandleFunc("/cart", handlers.ShowCart)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
