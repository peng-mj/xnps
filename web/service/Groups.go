package service

import (
	"tunpx/pkg/models"
)

type Group struct {
	Base
}

func NewGroup(db *Base) *Group {
	c := &Group{}
	c.Service(db.Driver)
	return c
}
func (c *Group) Create(group *models.Group) error {
	// maybe check unique group name
	return c.Orm(models.Group{}).Create(group).Error
}
