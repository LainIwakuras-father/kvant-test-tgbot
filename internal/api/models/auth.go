package models

type LoginInput struct {
	username    string `json:"username" binding:"required,username"`
	
}

// // TokenResponse defines the structure of the JWT token response
// // swagger:model
// type TokenResponse struct {
// 	Token string `json:"token"`
// }