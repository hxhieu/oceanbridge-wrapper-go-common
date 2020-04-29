package models

// LoginRequest struct
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResult struct
type LoginResult struct {
	AccessToken string `json:"accessToken"`
	Expired     int64  `json:"expired"`
	Error       string `json:"error"`
}
