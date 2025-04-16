package helpers

// UserKey используется как ключ для хранения данных пользователя в контексте
type UserKey struct{}

// UserAuth содержит информацию о аутентифицированном пользователе
type UserAuth struct {
	Username string
	Email    string
}
