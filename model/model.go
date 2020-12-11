package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Host struct {
	Id             int
	Name           string
	Autoname       string
	IP             string
	SNMP_community string
	Status         string
}

type Client struct {
	Id        int
	Name      string
	GrafanaId string
}

func (cl Client) String() string {
	return fmt.Sprintf("Client: %v. <%d> ", cl.Name, cl.Id)
}

type configurator struct {
	Db *sql.DB
}

func (c configurator) Close() error {
	return c.Db.Close()
}

func (c configurator) String() string {
	return "configuratro"
}

func NewConfigurator(username, password, host, dbname string, port int) (*configurator, error) {
	c := new(configurator)
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable application_name=go_pollers",
			host, port, username, password, dbname))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	c.Db = db
	return c, nil
}

func (c configurator) Clients() []Client {
	sqlStatement := "SELECT id, name, grafana_id FROM client order by name"
	var clients []Client

	var cl Client
	rows, err := c.Db.Query(sqlStatement)
	if err != nil {
		log.Warn(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&cl.Id, &cl.Name, &cl.GrafanaId)
		clients = append(clients, cl)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(err)
	}
	return clients
}

func (c configurator) HostsByClient(clientId int) []Host {
	sqlStatement := `SELECT id, name, autoname, ip, snmp_community, status FROM host WHERE client_id=$1 order by autoname`
	var h Host
	var hosts []Host
	rows, err := c.Db.Query(sqlStatement, clientId)
	if err != nil {
		log.Warn(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&h.Id, &h.Name, &h.Autoname, &h.IP, &h.SNMP_community, &h.Status)
		hosts = append(hosts, h)
	}

	err = rows.Err()
	if err != nil {
		log.Warn(err)
	}
	return hosts
}
