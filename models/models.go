package models

type Product struct {
	Name      string `json:"name"`
	Inventory int    `json:"inventory"`
	Price     int    `json:"price"`
}

type Login struct {
	User     int `json:"user"`
	Password int `json:"password"`
}

type User struct {
	Name     string `json:"name"`
	Password int    `json:"password"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

type Walet struct {
	Amount int `json:"amount"`
}

type SelectProduct struct {
	ProductID int `json:"product_id"`
	Count     int `json:"count"`
}
