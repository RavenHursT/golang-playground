package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ravenhurst/golang-playground/consts"
)

func LogoutHandler(context echo.Context) (err error) {
	context.SetCookie(&http.Cookie{
		Name:     consts.AUTH_TOKEN_COOKIE_NAME,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	return context.NoContent(http.StatusOK)
}
