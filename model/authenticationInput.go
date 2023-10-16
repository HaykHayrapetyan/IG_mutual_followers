package model

type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	InstaId  string `json:"instaId"  binding:"required"`
}
