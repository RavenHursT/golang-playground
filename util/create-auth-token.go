package util

import (
	"time"

	"github.com/ravenhurst/golang-playground/consts"
	"github.com/ravenhurst/golang-playground/structs"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func CreateAuthToken(sessionID int, username string, expiry time.Time) string {
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: consts.PrivateKey}, &signerOpts)
	if err != nil {
		panic(err)
	}

	builder := jwt.Signed(signer)

	claims := structs.Claims{
		Claims: &jwt.Claims{
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			Expiry:   jwt.NewNumericDate(expiry.UTC()),
		},
		Username:  username,
		SessionID: sessionID,
	}

	rawJWT, err := builder.Claims(claims).CompactSerialize()
	if err != nil {
		panic(err)
	}

	publicKey := &consts.PrivateKey.PublicKey
	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, nil)
	if err != nil {
		panic(err)
	}

	object, err := encrypter.Encrypt([]byte(rawJWT))
	if err != nil {
		panic(err)
	}

	ecryptedJWT, err := object.CompactSerialize()
	if err != nil {
		panic(err)
	}

	return string(ecryptedJWT)
}
