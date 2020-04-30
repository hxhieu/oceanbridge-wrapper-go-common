package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

// DecodeToken to get the claims/values of a JWT
func DecodeToken(token string, signingKey string) (*jwt.Token, *jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	jwt, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	return jwt, &claims, nil
}

// GenerateToken a JWT token
func GenerateToken(claims map[string]interface{}, expireIn *time.Duration, signingKey string) (*string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}
	jwtauth.SetIssuedNow(jwtClaims)
	jwtauth.SetExpiryIn(jwtClaims, time.Duration(60*time.Minute))

	tokenAuth := jwtauth.New(jwt.SigningMethodHS256.Name, []byte(signingKey), nil)
	_, token, err := tokenAuth.Encode(jwtClaims)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
