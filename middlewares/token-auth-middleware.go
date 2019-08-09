package middlewares

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/ravenhurst/golang-playground/consts"
	"github.com/ravenhurst/golang-playground/structs"
)

func getJwtKey(token *jwt.Token) (interface{}, error) {
	return []byte(consts.JWT_SIGNING_KEY), nil
}

func TokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		authCookie, err := context.Cookie(consts.AUTH_TOKEN_COOKIE_NAME)
		if err != nil {
			return context.NoContent(http.StatusForbidden)
		}

		claims := new(structs.Claims)
		token, err := jwt.ParseWithClaims(
			authCookie.Value,
			claims,
			getJwtKey,
		)

		if !token.Valid {
			// Tell the browser to delete the invalid auth-cookie
			context.SetCookie(&http.Cookie{
				Name:     consts.AUTH_TOKEN_COOKIE_NAME,
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			})

			if validationErr, ok := err.(*jwt.ValidationError); ok {
				tokenExpired := validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0
				tokenSignatureInvalid := err == jwt.ErrSignatureInvalid
				if tokenExpired || tokenSignatureInvalid {
					return context.NoContent(http.StatusUnauthorized)
				}
				return context.NoContent(http.StatusBadRequest)
			}
		}

		if err := next(context); err != nil {
			context.Error(err)
		}
		return nil
	}
}
