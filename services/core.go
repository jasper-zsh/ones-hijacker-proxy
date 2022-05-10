package services

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"gorm.io/gorm"
)

type HandlerAwareService struct {
	handler *handlers.ONESRequestHandler
}

func NewHandlerAwareService(handler *handlers.ONESRequestHandler) *HandlerAwareService {
	s := &HandlerAwareService{
		handler: handler,
	}
	return s
}

type DatabaseAwareService struct {
	db *gorm.DB
}

func NewDatabaseAwareService(db *gorm.DB) *DatabaseAwareService {
	return &DatabaseAwareService{
		db: db,
	}
}
