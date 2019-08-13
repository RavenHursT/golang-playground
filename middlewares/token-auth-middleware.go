package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/ravenhurst/golang-playground/consts"
	"github.com/ravenhurst/golang-playground/structs"
)

func TokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		authCookie, err := context.Cookie(consts.AUTH_TOKEN_COOKIE_NAME)
		if err != nil {
			return context.NoContent(http.StatusForbidden)
		}

		parsedJWT, err := jwt.ParseSigned(authCookie.Value)
		if err != nil {
			return context.NoContent(http.StatusBadRequest)
		}

		claims := structs.Claims{}
		err = parsedJWT.Claims(&consts.PrivateKey.PublicKey, &claims)
		if err != nil {
			return context.NoContent(http.StatusUnauthorized)
		}

		if err := next(context); err != nil {
			context.Error(err)
		}
		return nil
	}
}
