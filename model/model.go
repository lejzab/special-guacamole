package model

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

func MakeConnectionString(db string,
	user string,
	password string,
	host string,
	port int,
) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, db)
}

var db *sql.DB
var err error
var wg sync.WaitGroup

func Clients() []Client {
	//sqlStatement := `SELECT id, name FROM client WHERE status='Włączony'`
	sqlStatement := `SELECT id, name FROM client`
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
