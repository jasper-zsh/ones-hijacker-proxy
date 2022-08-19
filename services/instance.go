package services

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type InstanceServiceDeps struct {
	dig.In

	DB      *gorm.DB
	Handler *handlers.ONESRequestHandler
}

type InstanceService struct {
	deps InstanceServiceDeps
}

func NewInstanceService(deps InstanceServiceDeps) *InstanceService {
	return &InstanceService{
		deps,
	}
}

func (s *InstanceService) ListInstances() ([]*models.Instance, error) {
	var instances []*models.Instance
	q := s.deps.DB.Find(&instances)
	if q.Error != nil {
		return nil, q.Error
	}

	return instances, nil
}

func (s *InstanceService) SaveInstance(instance *models.Instance) error {
	q := s.deps.DB.Save(instance)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (s *InstanceService) DeleteInstance(id uint) error {
	q := s.deps.DB.Delete(&models.Instance{}, id)
	if q.Error != nil {
		return q.Error
	}

	return nil
}

func (s *InstanceService) SelectInstance(id uint) error {
	var instance *models.Instance
	q := s.deps.DB.First(&instance, id)
	if q.Error != nil {
		return q.Error
	}

	originInstance := s.deps.Handler.Instance
	originBinding := s.deps.Handler.Binding
	s.deps.Handler.Binding = nil
	s.deps.Handler.Instance = instance
	err := s.deps.Handler.Login(nil)
	if err != nil {
		s.deps.Handler.Instance = originInstance
		s.deps.Handler.Binding = originBinding
		return err
	}
	return nil
}
