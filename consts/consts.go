package consts

import (
	"crypto/rand"
	"crypto/rsa"
)

const AUTH_COOKIE_EXPIRE_MINS = 5
const AUTH_TOKEN_COOKIE_NAME = "auth_token"
const JWT_SIGNING_KEY = "{MYSUPER_SECRET_KEY}"
const PORT = "port"

var PrivateKey *rsa.PrivateKey

func init() {
	var err error
	PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
}
