package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
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

		object, err := jose.ParseEncrypted(authCookie.Value)
		if err != nil {
			panic(err)
		}

		decryptedByteArray, err := object.Decrypt(consts.PrivateKey)
		if err != nil {
			panic(err)
		}

		decryptedJWT := string(decryptedByteArray)

		parsedJWT, err := jwt.ParseSigned(decryptedJWT)
		if err != nil {
			return context.NoContent(http.StatusBadRequest)
		}

		claims := structs.Claims{}
		err = parsedJWT.Claims(&consts.PrivateKey.PublicKey, &claims)
		if err != nil {
			return context.String(http.StatusUnauthorized, "Unauthorized: Could not access token claims")
		}

		tokenExpiry := *claims.Claims.Expiry
		now := time.Now()
		expired := *jwt.NewNumericDate(now) >= tokenExpiry
		if expired {
			context.SetCookie(&http.Cookie{
				Name:     consts.AUTH_TOKEN_COOKIE_NAME,
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			})
			return context.String(http.StatusUnauthorized, "UnAuthorized: Auth token expired.")
		}

		if err := next(context); err != nil {
			context.Error(err)
		}
		return nil
	}
}
