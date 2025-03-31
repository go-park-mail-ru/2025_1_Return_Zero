package repository

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        int64 `sql:"id"`
	Username  string `sql:"username"`
	Password  string `sql:"password_hash"`
	Email     string `sql:"email"`
	Thubmnail string `sql:"thumbnail_url"`
}