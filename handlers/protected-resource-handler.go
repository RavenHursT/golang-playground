package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ravenhurst/golang-playground/structs"
)

func ProtectedResourceHandler(context echo.Context) (err error) {
	return context.JSON(http.StatusOK, structs.ProtectedResource{
		Foo: "FOO",
	})
}
