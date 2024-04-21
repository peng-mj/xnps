package service

import (
	"xnps/pkg/database"
)

type Base struct {
	*database.Driver
	//kv  *cache.Cache
}

func (b *Base) Service(db *database.Driver) *Base {
	b.Driver = db
	return b
}
