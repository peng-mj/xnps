package WebApi

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAllClient(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func GetClientByCondition(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func AddClient(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func DelClient(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
func EditClient(c echo.Context) (err error) {

	return c.String(http.StatusOK, "")
}
