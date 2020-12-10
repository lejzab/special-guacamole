package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
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
		fmt.Println(err)
		return
	}
	log.Info(config)
	//
	//db, err := sql.Open("postgres", model.MakeConnectionString(config.Database, config.Username, config.Password, config.Host, config.Port))
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//
	//err = db.Ping()
	//if err != nil {
	//	panic(err)
	//}
	//log.Info("Successfully connected!")
	//clients := model.Clients()
	////var hosts []Host
	//for _, c := range clients {
	//	log.Info(c.Name)
	//	//hosts = HostsByClient(c)
	//	//for idx, h := range hosts {
	//	//	fmt.Printf("Client %s, host no %d, name=[%s], password=[%s]\n", c.Name, idx, h.Name, h.SNMP_community)
	//	//}
	//}
	log.Info("STOP")
}
