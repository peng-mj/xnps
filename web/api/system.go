package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xnps/web/service"
)

/**********        USER          *********/

type System struct {
	kit *service.Base
}

func NewSystem(kit *service.Base) *System {
	return &System{kit: kit}
}

func (s *System) GetConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")

}
func (s *System) Update(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

// Init to write config to file,and remove temp config file
func (s *System) Init(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

// StaticInit  to load system init html and other static files
func (s *System) StaticInit(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func (s *System) StaticSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
