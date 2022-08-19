package dao

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("hijacker.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&models.Instance{},
		&models.Account{},
		&models.Binding{},
		&models.Variable{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
