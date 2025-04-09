package usecase

import (
	"io"
)

type User struct {
	ID        int64
	Email     string
	Username  string
	Avatar    io.Reader
	Password  string 
	AvatarUrl string
}

type ChangeUserData struct {
	Username    string 
	Email       string 
	Password    string 

	NewUsername string 
	NewEmail    string 
	NewPassword string 
}