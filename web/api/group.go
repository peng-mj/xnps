package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tunpx/web/service"
)

type Group struct {
	kit *service.Base
}

func NewGroup(dr *service.Base) *Group {
	return &Group{kit: dr}
}

func (b *Group) GetAll(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Group) GetByFilter(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Group) Create(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Group) Delete(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Group) Update(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
