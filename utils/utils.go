package utils

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func ReadConfiguration(filename string) (Configuration, error) {
	file, err := os.Open(filename)
	defer file.Close()
	configuration := Configuration{}
	if err != nil {
		return configuration, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, err
	} else {
		return configuration, nil
	}
}
