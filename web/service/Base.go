package service

import (
	"tunpx/pkg/database"
)

type Base struct {
	*database.Driver
	// kv  *cache.Cache
}

func (b *Base) Service(db *database.Driver) *Base {
	b.Driver = db
	return b
}
func NewBase(db *database.Driver) *Base {
	return &Base{Driver: db}
}
