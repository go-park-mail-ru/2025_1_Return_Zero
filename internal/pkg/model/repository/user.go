package repository

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        int64  `sql:"id"`
	Username  string `sql:"username"`
	Password  string `sql:"password_hash"`
	Email     string `sql:"email"`
	Thumbnail string `sql:"thumbnail_url"`
}

type ChangeUserData struct {
	Username    string `sql:"username"`
	Email       string `sql:"email"`
	Password    string `sql:"password_hash"`

	NewUsername string `sql:"username"`
	NewEmail    string `sql:"email"`
	NewPassword string `sql:"password_hash"`
}