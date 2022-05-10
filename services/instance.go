package services

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"gorm.io/gorm"
)

type InstanceService struct {
	*HandlerAwareService
	*DatabaseAwareService
}

func NewInstanceService(db *gorm.DB, handler *handlers.ONESRequestHandler) *InstanceService {
	return &InstanceService{
		NewHandlerAwareService(handler),
		NewDatabaseAwareService(db),
	}
}

func (s *InstanceService) ListInstances() ([]*models.Instance, error) {
	var instances []*models.Instance
	q := s.db.Find(&instances)
	if q.Error != nil {
		return nil, q.Error
	}

	return instances, nil
}

func (s *InstanceService) SaveInstance(instance *models.Instance) error {
	q := s.db.Save(instance)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (s *InstanceService) DeleteInstance(id uint) error {
	q := s.db.Delete(&models.Instance{}, id)
	if q.Error != nil {
		return q.Error
	}

	return nil
}

func (s *InstanceService) SelectInstance(id uint) error {
	var instance *models.Instance
	q := s.db.First(&instance, id)
	if q.Error != nil {
		return q.Error
	}

	originInstance := s.handler.Instance()
	originAuth := s.handler.AuthInfo()
	s.handler.ClearAuthInfo()
	s.handler.SetInstance(instance)
	err := s.handler.Login(nil)
	if err != nil {
		s.handler.SetInstance(originInstance)
		s.handler.SetAuthInfo(originAuth)
		return err
	}
	return nil
}
