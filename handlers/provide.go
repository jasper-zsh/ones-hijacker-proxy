package handlers

import "go.uber.org/dig"

func Provide(c *dig.Container) {
	err := c.Provide(NewONESRequestHandler)
	if err != nil {
		panic(err)
	}
}
