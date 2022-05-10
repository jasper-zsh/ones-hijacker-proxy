package services

import "go.uber.org/dig"

func Provide(c *dig.Container) {
	err := c.Provide(func(deps AccountServiceDeps) *AccountService {
		return NewAccountService(deps)
	})
	if err != nil {
		panic(err)
	}
	err = c.Provide(func(deps InstanceServiceDeps) *InstanceService {
		return NewInstanceService(deps)
	})
	if err != nil {
		panic(err)
	}
}
