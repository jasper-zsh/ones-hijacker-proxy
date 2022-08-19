package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Email      string `json:"email"`
	Password   string `json:"password"`
	Note       string `json:"note"`
	InstanceID uint
	Instance   *Instance
}
