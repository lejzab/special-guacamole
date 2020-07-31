package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "name of config file")
	flag.StringVar(&configFile, "c", "config.json", "name of config file")
	flag.Parse()
	config, err := readConfiguration(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err = sql.Open("postgres", MakeConnectionString(config))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Print("Successfully connected!")
	clients := Clients()
	var hosts []Host
	for _, c := range clients {
		hosts = HostsByClient(c)
		for idx, h := range hosts {
			fmt.Printf("Client %s, host no %d, name=[%s], password=[%s]\n", c.Name, idx, h.Name, h.SNMP_community)
		}
	}
}

func readConfiguration(filename string) (Configuration, error) {
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
