package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xnps/web/service"
)

type Tunnel struct {
	kit *service.Base
}

func NewTunnel(dr *service.Base) *Tunnel {
	return &Tunnel{kit: dr}
}

func (t *Tunnel) GetAll(c *gin.Context) {

	c.JSON(http.StatusOK, "")
}
func (t *Tunnel) GetFilter(c *gin.Context) {

	c.JSON(http.StatusOK, "")
}
func (t *Tunnel) GetByIds(c *gin.Context) {

	c.JSON(http.StatusOK, "")
}
func (t *Tunnel) Create(c *gin.Context) {

	c.JSON(http.StatusOK, "")
}
func (t *Tunnel) Delete(c *gin.Context) {

	c.JSON(http.StatusOK, "")
}

func (t *Tunnel) Update(c *gin.Context) {

	c.JSON(http.StatusOK, "")
}
