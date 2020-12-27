package jwthelper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var (
	//jwtSignatureKey is a secret key to hash the JWT Token
	jwtSignatureKey string
)

//CustomClaims - Represent object of claims. Encourage all claims is referred to this struct
type CustomClaims struct {
	jwt.StandardClaims
	Id        int    `json:"jti,omitempty"`
	Role      string `json:"role,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

func init() {
	jwtSignatureKey = os.Getenv("jwt.secretKey")
}

//NewWithClaims will return token with custom claims
func NewWithClaims(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSignatureKey))

	if err != nil {
		return "", err
	}
	return ss, nil
}

//VerifyTokenWithClaims will verify the validity of token and return the claims
func VerifyTokenWithClaims(token string) (*CustomClaims, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("error retrieving claims")
	}

	timeNow := time.Now().Unix()
	if claims.ExpiresAt < timeNow {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
