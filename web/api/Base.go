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

func GetCode(ctx *gin.Context) int {
	if v, ok := ctx.Get("errorCode"); ok {
		return v.(int)
	}
	return 0
}

func Replay(ctx *gin.Context, err error, data interface{}) (ReData dto.Response) {
	GetCode(ctx)
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
