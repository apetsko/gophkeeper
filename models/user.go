package models

type User struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserEntry struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	PasswordHash string  `json:"password_hash"`
	Balance      float64 `json:"balance"`
}
