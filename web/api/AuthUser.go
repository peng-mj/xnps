package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xnps/web/service"
)

type AuthUser struct {
	kit *service.Base
}

func NewUser(dr *service.Base) *AuthUser {
	return &AuthUser{kit: dr}
}

func (a *AuthUser) GetAllUser(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (a *AuthUser) GetUserByUid(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
func (a *AuthUser) Login(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}

func (a *AuthUser) UpdateUser(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
