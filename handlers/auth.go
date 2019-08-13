package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/ravenhurst/golang-playground/consts"
	"github.com/ravenhurst/golang-playground/structs"
	"github.com/ravenhurst/golang-playground/util"
)

var users = map[string]string{
	"foo": "fooPass",
	"bar": "barPass",
}

func getSessionIDFromSomePersistance() (n int, err error) {
	return rand.Intn(1000000), nil
}

func setTokenCookie(context echo.Context, token string, expiry time.Time) {
	authCookie := new(http.Cookie)
	authCookie.Name = consts.AUTH_TOKEN_COOKIE_NAME
	authCookie.Value = token // token.SignedString([]byte(consts.JWT_SIGNING_KEY))
	authCookie.Expires = expiry
	authCookie.HttpOnly = true
	context.SetCookie(authCookie)
}

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
	setTokenCookie(
		context, 
		util.CreateAuthToken(sessionID, creds.Username, expiry), 
		expiry,
	)

	return context.NoContent(http.StatusCreated)
}
