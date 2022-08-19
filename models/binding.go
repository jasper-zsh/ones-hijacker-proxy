package models

import "time"

type Binding struct {
	AccountID  uint      `json:"account_id" gorm:"primaryKey;autoIncrement:false"`
	InstanceID uint      `json:"instance_id" gorm:"primaryKey;autoIncrement:false"`
	UserUUID   string    `json:"user_uuid"`
	Token      string    `json:"token"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
