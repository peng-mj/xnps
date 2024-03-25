package WebApi

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"net/http"
)

type ReObj struct {
	Status string      `json:"Status"`
	Data   interface{} `json:"Data"`
}

func (r *ReObj) Error() string {
	return ""
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
	if marshal, err := json.Marshal(ReData); err == nil {

		return string(marshal)
	} else {
		slog.Error(err.Error())
	}
	return `{"status":ERROR","data":"Unexpected bugs, please file an issue with our github"}`
}
func Replay(err error, data interface{}) (ReData ReObj) {
	if err != nil {
		ReData.Status = "ERROR"
		ReData.Data = err.Error()
	} else {
		ReData.Status = "OK"
		ReData.Data = data
	}
	return
	//return `{"status":ERROR","data":"Unexpected bugs, please file an issue with our github"}`
}

func GetSystemStatus(c echo.Context) (err error) {
	return c.String(http.StatusOK, "")
}
func GetConnectInfo(c echo.Context) (err error) {
	return c.String(http.StatusOK, "")
}
