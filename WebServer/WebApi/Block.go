package WebApi

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAllBlockList(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func GetBlockListByCondition(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func AddFBlockList(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func DelBlockList(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func EditBlockList(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
