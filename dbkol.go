package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
)

type Host struct {
	ID             int
	Name           string
	Autoname       string
	IP             string
	SNMP_community string
	Status         string
}

type Client struct {
	ID   int
	Name string
}

func MakeConnectionString(config Configuration) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)
}

var db *sql.DB
var err error
var wg sync.WaitGroup

//func main() {
//	config, err := readConfiguration("config.json")
//	//fmt.Println(config, err)
//	psqlInfo := makeConnectionString(config)
//	db, err = sql.Open("postgres", psqlInfo)
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	err = db.Ping()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("Successfully connected!")
//
//	clients := clients()
//	var hosts []Host
//	for _, c := range clients {
//		hosts = hosts_by_client(c)
//		fmt.Printf("Client: %s, hosts count: %d\n", c.Name, len(hosts))
//	}
//	fmt.Println(len(clients), cap(clients))
//	fmt.Println("KONIEC")
//
//}

func Clients() []Client {
	sqlStatement := `SELECT id, name FROM client WHERE status='Włączony'`
	clients := make([]Client, 0, 8)

	var c Client
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&c.ID, &c.Name)
		clients = append(clients, c)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return clients
}

func HostsByClient(client Client) []Host {
	sqlStatement := `SELECT id, name, autoname, ip, snmp_community, status FROM host WHERE client_id=$1;`
	var h Host
	hosts := make([]Host, 0, 64)
	rows, err := db.Query(sqlStatement, client.ID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&h.ID, &h.Name, &h.Autoname, &h.IP, &h.SNMP_community, &h.Status)
		hosts = append(hosts, h)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return hosts
}
