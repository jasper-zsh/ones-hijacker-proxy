package models

import "gorm.io/gorm"

type Instance struct {
	gorm.Model
	BaseURL  string     `json:"base_url"`
	Note     string     `json:"note"`
	Accounts []*Account `json:"accounts"`
}
