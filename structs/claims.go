package structs

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	SessionID int `json:"sessionId"`
	jwt.StandardClaims
}