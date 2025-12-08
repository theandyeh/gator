package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSetUser(t *testing.T) {
	cfg := &Config{
		Db_url:          "postgresql://test",
		Current_db_user: "",
	}

	err := cfg.SetUser("newuser")
	if err != nil {
		t.Errorf("SetUser failed: %v", err)
	}

	if cfg.Current_db_user != "newuser" {
		t.Errorf("Expected user 'newuser', got '%s'", cfg.Current_db_user)
	}
}

func TestSetUserNoDbUrl(t *testing.T) {
	cfg := &Config{
		Db_url:          "",
		Current_db_user: "",
	}

	err := cfg.SetUser("newuser")
	if err == nil {
		t.Error("Expected error when db_url is empty")
	}
}

func TestWriteAndReadConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cfg := &Config{
		Db_url:          "postgresql://localhost:5432/test",
		Current_db_user: "testuser",
	}

	err := cfg.WriteConfig()
	if err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	configPath := filepath.Join(tempDir, Config_file_name)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	readCfg, err := Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if readCfg.Db_url != cfg.Db_url {
		t.Errorf("Expected db_url '%s', got '%s'", cfg.Db_url, readCfg.Db_url)
	}
	if readCfg.Current_db_user != cfg.Current_db_user {
		t.Errorf("Expected user '%s', got '%s'", cfg.Current_db_user, readCfg.Current_db_user)
	}
}

func TestReadInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	configPath := filepath.Join(tempDir, Config_file_name)
	err := os.WriteFile(configPath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err = Read()
	if err == nil {
		t.Error("Expected error when reading invalid JSON")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}

	if filepath.Base(path) != Config_file_name {
		t.Errorf("Expected filename '%s', got '%s'", Config_file_name, filepath.Base(path))
	}
}

func TestConfigJSONMarshaling(t *testing.T) {
	cfg := &Config{
		Db_url:          "postgresql://test",
		Current_db_user: "user123",
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled Config
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.Db_url != cfg.Db_url {
		t.Errorf("Expected db_url '%s', got '%s'", cfg.Db_url, unmarshaled.Db_url)
	}
	if unmarshaled.Current_db_user != cfg.Current_db_user {
		t.Errorf("Expected user '%s', got '%s'", cfg.Current_db_user, unmarshaled.Current_db_user)
	}
}
