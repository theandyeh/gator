package cmd

import (
	"fmt"

	"github.com/theandyeh/gator/internal/app"
)

func HandlerLogin(s *app.State, c Command) error {
	if len(c.args) < 1 {
		return fmt.Errorf("login handler error: no username provided for login command")
	}

	if err := s.Cfg.SetUser(c.args[0]); err != nil {
		return fmt.Errorf("login handler error: %w", err)
	}

	fmt.Printf("Successfuly set DB User: %s", c.args[0])
	return nil
}
