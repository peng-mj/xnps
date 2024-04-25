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
	auth := GetUser(ctx)
	if auth == nil || (auth != nil && auth.AuthLevel > 0) {
		RepError(ctx, http.StatusForbidden)
		return
	}
	user := service.NewAuthUser(a.kit).GetAllUser()
	Response(ctx, user)
}
func (a *AuthUser) GetUserByUid(ctx *gin.Context) {
	auth := GetUser(ctx)
	if auth == nil {
		RepError(ctx, http.StatusForbidden)
		return
	}
	user, err := service.NewAuthUser(a.kit).GetUserByUid(auth.Uid)
	if err != nil {
		RepError(ctx, http.StatusNotFound)
		return
	}
	Response(ctx, user)
}

func (a *AuthUser) UpdateUser(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}

func (a *AuthUser) CreateNewUser(ctx *gin.Context) {
	user := GetUser(ctx)
	if user == nil {

	}

	ctx.JSON(http.StatusOK, "")
}
func (a *AuthUser) CreateRootUser(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
