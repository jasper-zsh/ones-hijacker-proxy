package services

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/types"
	"gorm.io/gorm"
)

type AccountService struct {
	*HandlerAwareService
	*DatabaseAwareService
}

func NewAccountService(db *gorm.DB, handler *handlers.ONESRequestHandler) *AccountService {
	return &AccountService{
		NewHandlerAwareService(handler),
		NewDatabaseAwareService(db),
	}
}

func (s *AccountService) ListAccounts() ([]*models.Account, error) {
	var accounts []*models.Account
	q := s.db.Find(&accounts)
	if q.Error != nil {
		return nil, q.Error
	}
	return accounts, nil
}

func (s *AccountService) SaveAccount(account *models.Account) error {
	q := s.db.Save(account)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (s *AccountService) DeleteAccount(id uint) error {
	q := s.db.Delete(&models.Account{}, id)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (s *AccountService) SelectAccount(id uint) error {
	var account *models.Account
	q := s.db.First(&account, id)
	if q.Error != nil {
		return q.Error
	}

	originAccount := s.handler.Account()
	originAuth := s.handler.AuthInfo()
	s.handler.ClearAuthInfo()
	s.handler.SetAccount(account)
	s.handler.SetAuthUpdatedCallback(func(info *types.User) {
		account.Token = info.Token
		account.UserUUID = info.UUID
		s.db.Save(account)
	})
	err := s.handler.Login(nil)
	if err != nil {
		s.handler.SetAccount(originAccount)
		s.handler.SetAuthInfo(originAuth)
		return err
	}
	return nil
}
