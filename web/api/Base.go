package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"xnps/web/dto"
)

type ErCode int

func (e ErCode) ToMsg() string {
	return http.StatusText(int(e))
}

func Response(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, dto.Response{Code: http.StatusOK, Data: data})
	return
}
func RepError(ctx *gin.Context, code int) {
	ctx.JSON(code, dto.Response{Code: code, ErMsg: ErCode(code).ToMsg()})
	return
}

func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, strconv.FormatInt(time.Now().Unix(), 10))
	ctx.Abort()
}

func GetUser(ctx *gin.Context) *dto.User {
	if user, ok := ctx.Get("user"); ok {
		u := user.(dto.User)
		return &u
	}
	return nil
}
