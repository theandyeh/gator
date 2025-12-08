package cmd

import (
	"errors"

	"github.com/theandyeh/gator/internal/app"
)

type Command struct {
	name string
	args []string
}

type Commands struct {
	list map[string]func(*app.State, Command) error
}

func (c *Commands) Run(s *app.State, cm Command) error {
	if cmdFunc, exists := c.list[cm.name]; exists {
		return cmdFunc(s, cm)
	}
	return errors.New("cmd error: command not found")
}

func (c *Commands) Register(name string, f func(*app.State, Command) error) error {
	if _, exists := c.list[name]; exists {
		return errors.New("cmd error: command already registered")
	}
	c.list[name] = f
	return nil
}
