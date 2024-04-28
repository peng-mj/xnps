package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"tunpx/web/dto"
)

func Response(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, dto.Response{Code: http.StatusOK, Data: data})
}
func RepError(ctx *gin.Context, code int) {
	ctx.JSON(code, dto.Response{Code: code, ErMsg: dto.RspCode(code).String()})
}
func RepErrorWithMsg(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, dto.Response{Code: code, ErMsg: msg})
}

func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, strconv.FormatInt(time.Now().Unix(), 10))
}

func GetUser(ctx *gin.Context) *dto.User {
	if user, ok := ctx.Get("user"); ok {
		u := user.(dto.User)
		return &u
	}
	return nil
}
