package control

import "go.uber.org/dig"

func Provide(c *dig.Container) {
	err := c.Provide(func(deps ControlDeps) *Control {
		return NewControl(deps)
	})
	if err != nil {
		panic(err)
	}
}
