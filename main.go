package main

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"math/rand"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"

	"github.com/ravenhurst/golang-playground/foo"
	"github.com/ravenhurst/golang-playground/structs"
)

var users = map[string]string{
	"foo": "fooPass",
	"bar": "barPass",
}

const AUTH_COOKIE_EXPIRE_MINS = 5
const AUTH_TOKEN_COOKIE_NAME = "auth_token"
const JWT_SIGNING_KEY = "{MYSUPER_SECRET_KEY}"

type claimsStruct struct {
	Username string `json:"username"`
	SessionID int `json:"sessionId"`
	jwt.StandardClaims
}

func getSessionIDFromSomePersistance() (n int, err error) {
	return rand.Intn(1000000), nil
}

func authMiddleware(context echo.Context) (err error) {
	creds := new(structs.Credentials)
	if err = context.Bind(creds); err != nil {
		return context.JSON(http.StatusBadRequest, new(structs.ErrorResponseBody))
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		context.NoContent(http.StatusUnauthorized)
		return
	}
	
	var sessionID int
	if sessionID, err = getSessionIDFromSomePersistance(); err != nil {
		errBody := structs.ErrorResponseBody{
			Message: "This error might be what came back from persistance",
		}
		context.JSON(http.StatusInternalServerError, errBody)
	}

	expiry := time.Now().Add(AUTH_COOKIE_EXPIRE_MINS * time.Minute)

	claims := claimsStruct{
		Username: creds.Username,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expiry.Unix(),
		},
	}

	claimsJSON, _ := json.Marshal(claims)
	fmt.Println(string(claimsJSON))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	authCookie := new(http.Cookie)
	authCookie.Name = AUTH_TOKEN_COOKIE_NAME
	authCookie.Value, _ = token.SignedString([]byte(JWT_SIGNING_KEY))
	authCookie.Expires = expiry
	authCookie.HttpOnly = true
	context.SetCookie(authCookie)

	return context.JSON(http.StatusCreated, creds)
}

type protectedResource struct {
	SessionID int `json:"sessionId"`
}

func getJwtKey (token *jwt.Token) (interface{}, error){
	return []byte(JWT_SIGNING_KEY), nil
}

func protectedResourceMiddleware(context echo.Context) (err error) {
	authCookie, err := context.Cookie(AUTH_TOKEN_COOKIE_NAME)
	if err != nil {
		return context.NoContent(http.StatusForbidden)
	}
	
	claims := new(claimsStruct)
	token, err := jwt.ParseWithClaims(
		authCookie.Value,
		claims, 
		getJwtKey,
	)

	if !token.Valid {		
		// Tell the browser to delete the invalid auth-cookie
		context.SetCookie(&http.Cookie{
			Name: AUTH_TOKEN_COOKIE_NAME,
			Value: "",
			Expires: time.Unix(0, 0),
			HttpOnly: true,
		})

		if validationErr, ok := err.(*jwt.ValidationError); ok {
			tokenExpired := validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0
			tokenSignatureInvalid := err == jwt.ErrSignatureInvalid
			if ( tokenExpired || tokenSignatureInvalid ) {
				return context.NoContent(http.StatusUnauthorized)	
			}
			return context.NoContent(http.StatusBadRequest)
		}
	}

	var resource protectedResource
	resource.SessionID = claims.SessionID
	return context.JSON(http.StatusOK, resource)
}

func main() {
	fmt.Println("Hello World")
	fmt.Printf("foo => %s", foo.GetFoo())

	e := echo.New()
	e.POST("/auth", authMiddleware)
	e.GET("/protected-resource", protectedResourceMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
