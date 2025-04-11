package delivery

// UserToFront represents user data
// @Description User data
type UserToFront struct {
	ID       int64  `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RegisterData represents user registration data
// @Description User registration data
type RegisterData struct {
	Username string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,alphanum,stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
}

// LoginData represents user login data
// @Description User login data
type LoginData struct {
	Username string `json:"username" valid:"alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,stringlength(4|25)"`
	Email    string `json:"email" valid:"email,stringlength(5|30)"`
}
