package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tunpx/web/service"
)

type Block struct {
	kit *service.Base
}

func NewBlock(dr *service.Base) *Block {
	return &Block{kit: dr}
}

func (b *Block) GetAll(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Block) GetFilter(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Block) GetByIds(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Block) Create(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Block) Delete(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (b *Block) Update(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
