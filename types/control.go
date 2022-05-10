package types

import "github.com/jasper-zsh/ones-hijacker-proxy/models"

type AuthInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	User     *User  `json:"user"`
}

type StatusResponse struct {
	Account  *models.Account  `json:"account"`
	Instance *models.Instance `json:"instance"`
	ErrorMsg string           `json:"error_msg"`
}

type Timing struct {
	Period   string
	Duration int64
}
