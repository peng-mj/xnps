package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**********        USER          *********/

func GetSysConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")

}
func EditSYsConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
