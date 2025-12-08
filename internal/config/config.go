package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const Config_file_name = ".gatorconfig.json"

type Config struct {
	Db_url          string `json:"db_url"`
	Current_db_user string `json:"current_user_name"`
}

func Read() (Config, error) {
	config_file_path, err := GetConfigPath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.ReadFile(config_file_path)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config file: %w", err)
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	if len(c.Db_url) < 1 {
		return fmt.Errorf("DB URL is not set in config, cannot set user")
	}

	c.Current_db_user = username

	err := c.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) WriteConfig() error {
	config_file_path, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("Unable to get config file path to Write: %w", err)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Unable to marshal JSON to Write: %w", err)
	}

	if err := os.WriteFile(config_file_path, data, 0644); err != nil {
		return fmt.Errorf("Unable to write config file: %w", err)
	}

	return nil
}

func GetConfigPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error reading Home Directory: %w", err)
	}

	config_file_path := homedir + string(os.PathSeparator) + Config_file_name

	return config_file_path, nil
}
