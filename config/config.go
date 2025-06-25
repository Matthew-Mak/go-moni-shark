package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	Token     string
	AppId     string
	BotPrefix string

	config *Config
)

type Config struct {
	Token     string `json:"token"`
	AppId     string `json:"AppId"`
	BotPrefix string `json:"BotPrefix"`
}

// ReadConfig reads the config.json file and unmarshals it into the Config struct
func ReadConfig() error {
	log.Println("Reading config file")
	file, err := os.ReadFile("../config.json")

	if err != nil {
		return err
	}

	log.Println("Unmarshalling config file")

	// unmarhsall file into config struct
	err = json.Unmarshal(file, &config)

	if err != nil {
		log.Println("Error unmarshalling config file")
		return err
	}

	Token = config.Token
	AppId = config.AppId
	BotPrefix = config.BotPrefix

	log.Printf("Loaded BotPrefix = %q", BotPrefix)
	return nil
}
