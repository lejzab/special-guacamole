package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Db struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}
}

func (c Configuration) String() string {
	return fmt.Sprintf("conf4postgres: %v:****@%v:%v/%v", c.Db.Username, c.Db.Host, c.Db.Port, c.Db.Database)
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
