package structs

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

type Claims struct {
	*jwt.Claims
	Username string   `json:"username,omitempty"`
	SessionID int `json:"sessionId,omitEmpty"`
}