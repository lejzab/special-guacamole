package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"kolekcjoner/pollers/model"
	"kolekcjoner/pollers/utils"
)

func main() {
	log.Info("START")
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "name of config file")
	flag.StringVar(&configFile, "c", "config.json", "name of config file")
	flag.Parse()
	config, err := utils.ReadConfiguration(configFile)
	if err != nil {
		log.Warn(err)
		return
	}
	log.Debug(config)
	configurator, err := model.NewConfigurator(config.Db.Username, config.Db.Password, config.Db.Host, config.Db.Database, config.Db.Port)
	if err != nil {
		log.Warn(err)
		return
	}
	defer configurator.Close()

	for idx, c := range configurator.Clients() {
		log.Info(idx, c)
		hosts := configurator.HostsByClient(c.Id)
		log.Infof("Number of hosts: %d.", len(hosts))
	}
	log.Debug(configurator)
	log.Info("STOP")
}
