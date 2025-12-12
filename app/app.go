package app

import "shop/ent"

type App struct {
	DB  *ent.Client
	Key []byte
}
