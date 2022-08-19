package models

import "time"

type Variable struct {
	Key       string    `json:"key" gorm:"primaryKey"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
