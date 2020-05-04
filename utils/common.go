package utils

import (
	"crypto/sha512"
	"crypto/subtle"
)

// OceanBrigdeClaims struct
type OceanBrigdeClaims struct {
	Email  string
	Scopes string
}

// ClaimKeys defines the claim key constants
var ClaimKeys = OceanBrigdeClaims{
	Email:  "oceanbridge_email",
	Scopes: "oceanbridge_scopes",
}

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
// https://github.com/go-macaron/auth/blob/884c0e6c9b92ceebf12101d6cf1417fe0f61bcac/util.go
func SecureCompare(given string, actual string) bool {
	givenSha := sha512.Sum512([]byte(given))
	actualSha := sha512.Sum512([]byte(actual))

	return subtle.ConstantTimeCompare(givenSha[:], actualSha[:]) == 1
}
