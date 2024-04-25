package models

type User struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}
