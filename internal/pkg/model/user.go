package model

type User struct {
	ID       uint   `json:"-"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" valid:"email"`
}

type UserToFront struct {
	ID       uint   `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterData struct {
	Username string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,alphanum,stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
}

type LoginData struct {
	Username string `json:"username" valid:"alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,stringlength(4|25)"`
	Email    string `json:"email" valid:"email,stringlength(5|30)"`
}
