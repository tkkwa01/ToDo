package request

type CreateUser struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserLogin struct {
	Session  bool   `json:"session"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
