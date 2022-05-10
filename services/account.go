package services

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/types"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type AccountServiceDeps struct {
	dig.In

	DB      *gorm.DB
	Handler *handlers.ONESRequestHandler
}

type AccountService struct {
	deps AccountServiceDeps
}

func NewAccountService(deps AccountServiceDeps) *AccountService {
	return &AccountService{
		deps,
	}
}

func (s *AccountService) ListAccounts() ([]*models.Account, error) {
	var accounts []*models.Account
	q := s.deps.DB.Find(&accounts)
	if q.Error != nil {
		return nil, q.Error
	}
	return accounts, nil
}

func (s *AccountService) SaveAccount(account *models.Account) error {
	q := s.deps.DB.Save(account)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (s *AccountService) DeleteAccount(id uint) error {
	q := s.deps.DB.Delete(&models.Account{}, id)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (s *AccountService) SelectAccount(id uint) error {
	var account *models.Account
	q := s.deps.DB.First(&account, id)
	if q.Error != nil {
		return q.Error
	}

	originAccount := s.deps.Handler.Account()
	originAuth := s.deps.Handler.AuthInfo()
	s.deps.Handler.ClearAuthInfo()
	s.deps.Handler.SetAccount(account)
	s.deps.Handler.SetAuthUpdatedCallback(func(info *types.User) {
		account.Token = info.Token
		account.UserUUID = info.UUID
		s.deps.DB.Save(account)
	})
	err := s.deps.Handler.Login(nil)
	if err != nil {
		s.deps.Handler.SetAccount(originAccount)
		s.deps.Handler.SetAuthInfo(originAuth)
		return err
	}
	return nil
}
