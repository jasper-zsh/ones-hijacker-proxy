package dao

import (
	"go.uber.org/dig"
)

func Provide(c *dig.Container) {
	err := c.Provide(InitDB)
	if err != nil {
		panic(err)
	}
}
