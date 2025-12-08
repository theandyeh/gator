package cmd

import (
	"errors"
	"testing"

	"github.com/theandyeh/gator/internal/app"
	"github.com/theandyeh/gator/internal/config"
)

func TestCreateCommandsList(t *testing.T) {
	cmdList := CreateCommandsList()
	if cmdList == nil {
		t.Fatal("CreateCommandsList returned nil")
	}
	if cmdList.List == nil {
		t.Fatal("Commands.List map is nil")
	}
	if len(cmdList.List) != 0 {
		t.Errorf("Expected empty command list, got %d commands", len(cmdList.List))
	}
}

func TestRegister(t *testing.T) {
	cmdList := CreateCommandsList()
	mockHandler := func(s *app.State, c Command) error {
		return nil
	}

	err := cmdList.Register("test", mockHandler)
	if err != nil {
		t.Fatalf("Failed to register command: %v", err)
	}

	if len(cmdList.List) != 1 {
		t.Errorf("Expected 1 command, got %d", len(cmdList.List))
	}

	if _, exists := cmdList.List["test"]; !exists {
		t.Error("Registered command not found in list")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	cmdList := CreateCommandsList()
	mockHandler := func(s *app.State, c Command) error {
		return nil
	}

	cmdList.Register("test", mockHandler)
	err := cmdList.Register("test", mockHandler)

	if err == nil {
		t.Error("Expected error when registering duplicate command")
	}
}

func TestRun(t *testing.T) {
	cmdList := CreateCommandsList()
	called := false
	mockHandler := func(s *app.State, c Command) error {
		called = true
		return nil
	}

	cmdList.Register("test", mockHandler)
	state := &app.State{Cfg: &config.Config{}}
	cmd := Command{Name: "test", Args: []string{}}

	err := cmdList.Run(state, cmd)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestRunNonExistent(t *testing.T) {
	cmdList := CreateCommandsList()
	state := &app.State{Cfg: &config.Config{}}
	cmd := Command{Name: "nonexistent", Args: []string{}}

	err := cmdList.Run(state, cmd)
	if err == nil {
		t.Error("Expected error for non-existent command")
	}
}

func TestRunWithError(t *testing.T) {
	cmdList := CreateCommandsList()
	expectedErr := errors.New("handler error")
	mockHandler := func(s *app.State, c Command) error {
		return expectedErr
	}

	cmdList.Register("test", mockHandler)
	state := &app.State{Cfg: &config.Config{}}
	cmd := Command{Name: "test", Args: []string{}}

	err := cmdList.Run(state, cmd)
	if err == nil {
		t.Error("Expected error from handler")
	}
	if err != expectedErr {
		t.Errorf("Expected %v, got %v", expectedErr, err)
	}
}

func TestRunWithArgs(t *testing.T) {
	cmdList := CreateCommandsList()
	var receivedArgs []string
	mockHandler := func(s *app.State, c Command) error {
		receivedArgs = c.Args
		return nil
	}

	cmdList.Register("test", mockHandler)
	state := &app.State{Cfg: &config.Config{}}
	expectedArgs := []string{"arg1", "arg2"}
	cmd := Command{Name: "test", Args: expectedArgs}

	cmdList.Run(state, cmd)

	if len(receivedArgs) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(receivedArgs))
	}
	for i, arg := range expectedArgs {
		if receivedArgs[i] != arg {
			t.Errorf("Arg %d: expected %s, got %s", i, arg, receivedArgs[i])
		}
	}
}
