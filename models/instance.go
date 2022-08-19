package models

import "gorm.io/gorm"

type Instance struct {
	gorm.Model
	Mode     string     `json:"mode"`
	BaseURL  string     `json:"base_url"`
	Note     string     `json:"note"`
	Accounts []*Account `json:"accounts"`
}
