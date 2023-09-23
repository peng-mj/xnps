package WebApi

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReObj struct {
	Status string      `json:"Status"`
	Data   interface{} `json:"Data"`
}

func ReDara(err error, data interface{}) string {
	ReData := ReObj{}
	if err != nil {
		ReData.Status = "ERROR"
		ReData.Data = err.Error()
	} else {
		ReData.Status = "OK"
		ReData.Data = data
	}
	if marshal, err := json.Marshal(ReData); err != nil {
		return string(marshal)
	}
	return `{"status":ERROR","data":"Unexpected bugs, please file an issue with our github"}`
}

func GetSystemStatus(c echo.Context) (err error) {
	return c.String(http.StatusOK, "")
}
func GetConnectInfo(c echo.Context) (err error) {
	return c.String(http.StatusOK, "")
}

/**********        BLOCK_LIST          *********/
