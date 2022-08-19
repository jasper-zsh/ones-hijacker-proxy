package models

import "gorm.io/gorm"

type Instance struct {
	gorm.Model
	Mode     string     `json:"mode"`
	Project  string     `json:"project"`
	Wiki     string     `json:"wiki"`
	Note     string     `json:"note"`
	Accounts []*Account `json:"accounts"`
}
