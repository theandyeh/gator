package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/theandyeh/gator/internal/app"
	"github.com/theandyeh/gator/internal/database"
)

func HandlerLogin(s *app.State, c Command) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("login handler error: no username provided for login command")
	}

	if _, err := s.Db.GetUser(context.Background(), c.Args[0]); err != nil {
		return fmt.Errorf("login handler error: %w", err)
	}

	if err := s.Cfg.SetUser(c.Args[0]); err != nil {
		return fmt.Errorf("login handler error: %w", err)
	}

	fmt.Printf("Successfuly set DB User: %s", c.Args[0])
	return nil
}

func HandlerRegister(s *app.State, c Command) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("register handler error: no username provided for register command")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
	}

	if u, err := s.Db.GetUser(context.Background(), c.Args[0]); err == nil && u.Name == c.Args[0] {
		fmt.Printf("register handler error: user %s already exists\n", c.Args[0])
		os.Exit(1)
	} else if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return fmt.Errorf("register handler error: %w", err)
		}
	}

	if _, err := s.Db.CreateUser(context.Background(), user); err != nil {
		return fmt.Errorf("register handler error: %w", err)
	}

	s.Cfg.Current_db_user = c.Args[0]
	s.Cfg.SetUser(c.Args[0])

	fmt.Printf("Successfuly registered and set DB User: %s\n", c.Args[0])
	fmt.Printf("User details:\n %v", user)
	return nil
}

func HandlerReset(s *app.State, c Command) error {
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("reset handler error: %w", err)
	}

	s.Cfg.Current_db_user = ""
	s.Cfg.SetUser("")

	fmt.Println("Successfully reset database users and cleared current user in config")
	return nil
}

func HandlerUsers(s *app.State, c Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("users handler error: %w", err)
	}

	fmt.Println("Registered Users:")
	for _, user := range users {
		if user.Name == s.Cfg.Current_db_user {
			fmt.Printf("- %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("- %s\n", user.Name)
	}

	return nil
}
