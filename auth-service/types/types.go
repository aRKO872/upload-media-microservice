package types

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id			 string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Picture struct {
	UserId	 string `json:"user_id"`
	PictureURL string `json:"picture_url"`
}

type LoginUserResponse struct {
	Status bool `json:"status"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ErrMsg string `json:"error_message"`
	ErrDetails []string `json:"error_details"`
}

type UploadPicResponse struct {
	Status bool `json:"status"`
	PictureURL string `json:"picture_url"`
	ErrMsg string `json:"error_message"`
}

type Claims struct {
	UserId string 
	jwt.RegisteredClaims
}

type ProducedOrConsumedMessage struct {
	Email string					`json:"email"`
	PictureURL string			`json:"picture_url"`
}