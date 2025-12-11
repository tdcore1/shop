package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"shop/auth"
	"shop/db"
	"shop/ent/cart"
	"shop/ent/walet"
	"shop/models"
)

// Login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login models.Login
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &login)

	u, _ := db.Client.User.Get(context.Background(), login.User)
	if u.Password != login.Password {
		fmt.Fprintln(w, "false")
		return
	}

	token := auth.CreateToken(login.User)
	w.Write([]byte(token))
}

// Product
func HandlerProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintln(w, "method false")
		return
	}
	result, _ := db.Client.Product.Query().All(context.Background())
	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func HandlerAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(w, "method false")
		return
	}
	var add models.Product
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &add)

	ad, err := db.Client.Product.Create().
		SetName(add.Name).
		SetInventory(add.Inventory).
		SetPrice(add.Price).
		Save(context.Background())
	fmt.Println(ad)
	if err == nil {
		fmt.Fprintln(w, "done")
	}
}

// User
func HandlerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(w, "method false")
		return
	}
	var user models.User
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &user)

	userdb, err := db.Client.User.Create().
		SetName(user.Name).
		SetLastName(user.LastName).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(context.Background())
	fmt.Println(userdb)
	if err == nil {
		fmt.Fprintln(w, "done")
	}
}

// Walet
func HandlerWalet(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.Check(w, r)
	if !ok {
		return
	}
	var walet models.Walet
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &walet)

	waletdb, err := db.Client.Walet.Create().
		SetAmount(walet.Amount).
		SetUserID(userID).Save(context.Background())
	if err == nil {
		fmt.Fprintln(w, "done")
		fmt.Fprintln(w, waletdb)
	}
}

// Cart
func CartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.Check(w, r)
	if !ok {
		return
	}
	var sel models.SelectProduct
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &sel)

	ctx := context.Background()
	productSelect, _ := db.Client.Product.Get(ctx, sel.ProductID)
	total := sel.Count * productSelect.Price

	wal, _ := db.Client.Walet.Query().
		Where(walet.UserID(userID)).
		Only(ctx)

	if wal.Amount < total {
		fmt.Fprintln(w, "not enough money")
		return
	}

	if productSelect.Inventory < sel.Count {
		fmt.Fprintln(w, "not enough inventory")
		db.Client.Product.DeleteOneID(productSelect.ID).Exec(ctx)
		return
	}

	wal.Update().SetAmount(wal.Amount - total).Save(ctx)
	productSelect.Update().SetInventory(productSelect.Inventory - sel.Count).Save(ctx)

	record, err := db.Client.Cart.Query().
		Where(cart.ProductID(sel.ProductID), cart.UserID(userID)).
		Only(ctx)
	if err == nil {
		record.Update().SetCount(record.Count + sel.Count).Save(ctx)
	} else {
		db.Client.Cart.Create().
			SetUserID(userID).SetCount(sel.Count).SetProductID(sel.ProductID).Save(ctx)
	}
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.Check(w, r)
	if !ok {
		return
	}
	items, _ := db.Client.Cart.Query().
		Where(cart.UserID(userID)).
		All(context.Background())

	data, _ := json.Marshal(items)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	wal, _ := db.Client.Walet.Query().
		Where(walet.UserID(userID)).
		Only(context.Background())

	fmt.Fprintln(w, "your amount:", wal.Amount)
}
