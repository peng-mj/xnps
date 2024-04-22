package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xnps/web/service"
)

type Client struct {
	kit *service.Base
}

func NewClient(dr *service.Base) *Client {
	return &Client{kit: dr}
}

func (c *Client) GetAll(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}
func (c *Client) GetFilter(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}
func (c *Client) GetByIds(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}

func (c *Client) Create(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}
func (c *Client) Delete(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}
func (c *Client) Update(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}
