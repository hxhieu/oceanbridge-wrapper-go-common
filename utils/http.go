package utils

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/hxhieu/oceanbridge-wrapper-go-common/models"
	"github.com/hxhieu/oceanbridge-wrapper-go-common/persistent"
	userstores "github.com/hxhieu/oceanbridge-wrapper-go-common/persistent/userstores"
)

// DevMode defines if the application is run on dev mode
var DevMode bool = strings.ToLower(os.Getenv("CARGOWISE_OCEANBRIDGE_REST_WRAPPER_PORT")) == "true"

// SetCors settings the CORS headers
func SetCors(w *http.ResponseWriter) {
	// TODO: Just allow all for the discovery stage
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// GetValidSession gets a valid session from the list, also returns the list of invalid sessions for clean up purpose
func GetValidSession(sessions *[]models.AuthSession, signingKey string) (*models.AuthSession, []string) {
	var invalidSessions []string
	var foundSession *models.AuthSession

	// Validate all sessions
	for _, session := range *sessions {
		_, claims, err := DecodeToken(session.AccessToken, signingKey)
		// Invalid token
		if err != nil {
			invalidSessions = append(invalidSessions, session.ID)
			continue
		}
		// Invalid claims
		if err = claims.Valid(); err != nil {
			invalidSessions = append(invalidSessions, session.ID)
			continue
		}

		foundSession = &session
		break
	}

	return foundSession, invalidSessions
}

// ParseAccessToken will try to parse the JWT,
// from the Authorization header then return the user instance as well as the decoded claims map
func ParseAccessToken(r *http.Request, w *http.ResponseWriter) (*persistent.User, *jwt.MapClaims) {
	// Just use a local settings instead of burning the storage quotas
	if DevMode {
		return &persistent.User{
			Email:    os.Getenv("CARGOWISE_OCEANBRIDGE_REST_WRAPPER_USER"),
			Password: os.Getenv("CARGOWISE_OCEANBRIDGE_REST_WRAPPER_PASSWORD"),
		}, nil
	}

	// Get get signing key
	secret := os.Getenv("CARGOWISE_OCEANBRIDGE_REST_WRAPPER_JWT_SECRET")
	if len(secret) <= 0 {
		http.Error(*w, "Server is not set up properly, missing JWT sign key.", http.StatusInternalServerError)
		return nil, nil
	}

	// Try parse auth header
	auth := r.Header.Get("Authorization")
	if len(auth) <= 0 {
		http.Error(*w, "Bearer access token is required.", http.StatusUnauthorized)
		return nil, nil
	}
	segs := strings.Fields(auth)
	if len(segs) != 2 {
		http.Error(*w, "Invalid authentication headers scheme, expecting a Bearer token.", http.StatusUnauthorized)
		return nil, nil
	}
	scheme := strings.ToUpper(segs[0])
	if !SecureCompare(scheme, "BEARER") {
		http.Error(*w, "Invalid authentication headers scheme, expecting a Bearer token.", http.StatusUnauthorized)
		return nil, nil
	}
	jwt := segs[1]
	_, claims, err := DecodeToken(jwt, secret)
	if err != nil {
		http.Error(*w, "Invalid access token.", http.StatusUnauthorized)
		return nil, nil
	}
	// Try load user sessions
	email := (*claims)[ClaimKeys.Email].(string)
	userStore, err := userstores.NewUserStoreFirestore()
	if err != nil {
		http.Error(*w, "An error occured while trying to access the database.", http.StatusInternalServerError)
		return nil, nil
	}
	sessions, user, err := userStore.GetSessions(email)
	// Get a valid session
	// Don't need to care about invalid sessions as they will get cleaned up as soon as re-login happens
	session, _ := GetValidSession(sessions, secret)
	if session == nil || !SecureCompare(jwt, session.AccessToken) {
		// Clean up the invalid session
		if session != nil {
			userStore.Logout(session.ID)
		}
		http.Error(*w, "Invalid or expired session.", http.StatusForbidden)
		return nil, nil
	}

	return user, claims
}
