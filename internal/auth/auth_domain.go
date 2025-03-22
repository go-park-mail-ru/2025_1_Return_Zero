package auth

type RegisterUserData struct {
	Username string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email"`
}

type LoginUserData struct {
	Username string `json:"username" valid:"alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,stringlength(4|25)"`
	Email    string `json:"email" valid:"email"`
}

type UserToFront struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
