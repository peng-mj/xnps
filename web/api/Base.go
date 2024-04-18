package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReObj struct {
	Code  int         `json:"code"`
	ErMsg string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type ErCode int

func (e ErCode) ToMsg() string {
	return http.StatusText(int(e))
}

func GetUser(ctx *gin.Context) string {
	if v, ok := ctx.Get("user"); ok {
		return v.(string)
	}
	return ""
}
func GetCode(ctx *gin.Context) int {
	if v, ok := ctx.Get("errorCode"); ok {
		return v.(int)
	}
	return 0
}

func Replay(ctx *gin.Context, err error, data interface{}) (ReData ReObj) {
	GetCode(ctx)
	return
}
