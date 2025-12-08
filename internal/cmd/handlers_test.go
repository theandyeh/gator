package cmd

import (
	"testing"

	"github.com/theandyeh/gator/internal/app"
	"github.com/theandyeh/gator/internal/config"
)

func TestHandlerLoginNoArgs(t *testing.T) {
	state := &app.State{
		Cfg: &config.Config{
			Db_url: "postgresql://test",
		},
	}
	cmd := Command{Name: "login", Args: []string{}}

	err := HandlerLogin(state, cmd)
	if err == nil {
		t.Error("Expected error when no username provided")
	}
}

func TestHandlerLoginNoDbUrl(t *testing.T) {
	state := &app.State{
		Cfg: &config.Config{},
	}
	cmd := Command{Name: "login", Args: []string{"testuser"}}

	err := HandlerLogin(state, cmd)
	if err == nil {
		t.Error("Expected error when db_url not set")
	}
}

func TestHandlerLoginSuccess(t *testing.T) {
	state := &app.State{
		Cfg: &config.Config{
			Db_url: "postgresql://test",
		},
	}
	cmd := Command{Name: "login", Args: []string{"testuser"}}

	err := HandlerLogin(state, cmd)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if state.Cfg.Current_db_user != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", state.Cfg.Current_db_user)
	}
}
