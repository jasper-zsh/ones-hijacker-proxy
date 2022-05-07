package types

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User *User `json:"user"`
}

type User struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Token string `json:"token"`
}
