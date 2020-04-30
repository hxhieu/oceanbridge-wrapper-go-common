package models

// LoginRequest struct
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResult struct
type LoginResult struct {
	AccessToken string `json:"accessToken"`
}

// AuthSession struct
type AuthSession struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}
