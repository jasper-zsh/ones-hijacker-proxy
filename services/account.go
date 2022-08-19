package services

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"strconv"
)

const (
	KeyDefaultAccountID = "DEFAULT_ACCOUNT_ID"
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

func (s *AccountService) SelectAccount(id uint) (err error) {
	defer func() {
		if err == nil {
			s.deps.DB.Save(&models.Variable{
				Key:   KeyDefaultAccountID,
				Value: strconv.Itoa(int(id)),
			})
		}
	}()
	var account *models.Account
	q := s.deps.DB.First(&account, id)
	if q.Error != nil {
		err = q.Error
		return
	}

	var binding *models.Binding
	q = s.deps.DB.Find(&binding, map[string]interface{}{
		"account_id":  account.ID,
		"instance_id": s.deps.Handler.Instance.ID,
	})
	if q.Error != nil {
		err = q.Error
		return
	}
	if binding != nil {
		s.deps.Handler.Account = account
		s.deps.Handler.Binding = binding
		return
	}

	originAccount := s.deps.Handler.Account
	originBinding := s.deps.Handler.Binding
	s.deps.Handler.Binding = nil
	s.deps.Handler.Account = account
	err = s.deps.Handler.Login(nil)
	if err != nil {
		s.deps.Handler.Account = originAccount
		s.deps.Handler.Binding = originBinding
		return
	}
	return
}
