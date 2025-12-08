package cmd

import (
	"errors"

	"github.com/theandyeh/gator/internal/app"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	List map[string]func(*app.State, Command) error
}

func (c *Commands) Run(s *app.State, cm Command) error {
	if cmdFunc, exists := c.List[cm.Name]; exists {
		return cmdFunc(s, cm)
	}
	return errors.New("cmd error: command not found")
}

func (c *Commands) Register(name string, f func(*app.State, Command) error) error {
	if _, exists := c.List[name]; exists {
		return errors.New("cmd error: command already registered")
	}
	c.List[name] = f
	return nil
}

func CreateCommandsList() *Commands {
	return &Commands{
		List: make(map[string]func(*app.State, Command) error),
	}
}
