package persistent

import "github.com/hxhieu/oceanbridge-wrapper-go-common/models"

// UserStore service
type UserStore interface {
	Login(email string, password string) (*[]models.AuthSession, error)
	Logout(sessionIDs ...string) error
	CreateSession(email string, accessToken string) error
	SetToken(sessionID string, accessToken string) error
	GetSession(sessionID string) (*models.AuthSession, *User, error)
}
