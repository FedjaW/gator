package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DbUrl string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}


func Read() (Config, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return Config{}, err
    }

    filepath := fmt.Sprintf("%s/%s", homeDir, configFileName)
    fileData, err := os.ReadFile(filepath)
    if err != nil {
        return Config{}, err
    }

    cfg := Config{}

    err = json.Unmarshal(fileData, &cfg)
    if err != nil {
        return Config{}, err
    }

    return cfg, nil
}

func (cfg *Config) SetUser(currentUserName string) error {
    cfg.CurrentUserName = currentUserName

    jsonData, err := json.Marshal(*cfg)
    if err != nil {
        return err
    }


    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    filepath := fmt.Sprintf("%s/%s", homeDir, configFileName)

    err = os.WriteFile(filepath, jsonData, 0644)
    if err != nil {
        return err
    }
    
    return nil
}
