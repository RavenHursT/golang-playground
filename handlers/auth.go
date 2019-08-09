package handlers

import (
	"net/http"
	"math/rand"
	"time"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	
	"github.com/ravenhurst/golang-playground/consts"
	"github.com/ravenhurst/golang-playground/structs"
)

var users = map[string]string{
	"foo": "fooPass",
	"bar": "barPass",
}

func getSessionIDFromSomePersistance() (n int, err error) {
	return rand.Intn(1000000), nil
}

// AuthHandler comment
func AuthHandler(context echo.Context) (err error) {
	creds := new(structs.Credentials)
	if err = context.Bind(creds); err != nil {
		return context.JSON(http.StatusBadRequest, new(structs.ErrorResponseBody))
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		return context.NoContent(http.StatusUnauthorized)
	}
	
	var sessionID int
	if sessionID, err = getSessionIDFromSomePersistance(); err != nil {
		errBody := structs.ErrorResponseBody{
			Message: "This error might be what came back from persistance",
		}
		return context.JSON(http.StatusInternalServerError, errBody)
	}

	expiry := time.Now().Add(consts.AUTH_COOKIE_EXPIRE_MINS * time.Minute)

	claims := structs.Claims{
		Username: creds.Username,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expiry.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	authCookie := new(http.Cookie)
	authCookie.Name = consts.AUTH_TOKEN_COOKIE_NAME
	authCookie.Value, _ = token.SignedString([]byte(consts.JWT_SIGNING_KEY))
	authCookie.Expires = expiry
	authCookie.HttpOnly = true
	context.SetCookie(authCookie)

	return context.NoContent(http.StatusCreated)
}