package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"shop/ent/cart"
	"shop/ent/walet"
	"shop/models"
)

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login models.Login

	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &login)

	u, err := a.DB.User.Get(context.Background(), login.User)
	if err != nil || u.Password != login.Password {
		http.Error(w, "wrong user or pass", 401)
		return
	}

	token, _ := a.CreateToken(u.ID)
	w.Write([]byte(token))
}

func (a *App) AddProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &p)

	a.DB.Product.Create().
		SetName(p.Name).
		SetInventory(p.Inventory).
		SetPrice(p.Price).
		Save(context.Background())

	fmt.Fprintln(w, "product added")
}

func (a *App) GetProducts(w http.ResponseWriter, r *http.Request) {
	result, _ := a.DB.Product.Query().All(context.Background())
	data, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (a *App) AddWalet(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.Check(w, r)
	if !ok {
		return
	}

	var wlt models.Walet
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &wlt)

	a.DB.Walet.Create().
		SetUserID(userID).
		SetAmount(wlt.Amount).
		Save(context.Background())

	fmt.Fprintln(w, "walet added")
}

func (a *App) SelectProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.Check(w, r)
	if !ok {
		return
	}

	var sp models.SelectProduct
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &sp)

	ctx := context.Background()

	product, _ := a.DB.Product.Get(ctx, sp.ProductID)
	total := product.Price * sp.Count

	waletDB, _ := a.DB.Walet.Query().Where(walet.UserID(userID)).Only(ctx)

	if waletDB.Amount < total {
		fmt.Fprintln(w, "not enough money")
		return
	}

	if product.Inventory < sp.Count {
		fmt.Fprintln(w, "not enough inventory")
		a.DB.Product.DeleteOneID(product.ID).Exec(ctx)
		return
	}

	product.Update().SetInventory(product.Inventory - sp.Count).Save(ctx)
	waletDB.Update().SetAmount(waletDB.Amount - total).Save(ctx)

	record, err := a.DB.Cart.Query().
		Where(cart.UserID(userID), cart.ProductID(sp.ProductID)).
		Only(ctx)

	if err == nil {
		record.Update().SetCount(record.Count + sp.Count).Save(ctx)
	} else {
		a.DB.Cart.Create().
			SetUserID(userID).
			SetProductID(sp.ProductID).
			SetCount(sp.Count).
			Save(ctx)
	}

	fmt.Fprintln(w, "added to cart")
}

func (a *App) ShowCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.Check(w, r)
	if !ok {
		return
	}

	records, _ := a.DB.Cart.Query().
		Where(cart.UserID(userID)).
		All(context.Background())

	data, _ := json.Marshal(records)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
