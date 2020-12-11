package model

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type ModelError struct {
	OriginalMessage string
	OriginalHint    string
	Message         string
}

func (e ModelError) Error() string {
	return fmt.Sprintf("ModelError: %s. Original Error: %s. %s", e.Message, e.OriginalMessage, e.OriginalHint)
}

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

type PingProfile struct {
	Id             int
	PacketCount    int
	PacketInterval int
	PacketSize     int
	Success        int
	Timeout        int
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
func (c configurator) PingProfiles() ([]PingProfile, error) {
	var profiles []PingProfile

	return profiles, nil
}
func (c configurator) Clients() ([]Client, error) {
	sqlStatement := "SELECT id, name, grafana_id FROM client order by name limit 2"
	var clients []Client
	var cl Client

	rows, err := c.Db.Query(sqlStatement)
	if err != nil {
		return nil, modelError("sql statemant error", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&cl.Id, &cl.Name, &cl.GrafanaId)
		if err != nil {
			return nil, modelError("error fetching clients", err)
		}
		clients = append(clients, cl)
	}
	err = rows.Err()
	if err != nil {
		return nil, modelError("error fetching clients", err)
	}
	return clients, nil
}

func (c configurator) HostsByClient(clientId int) ([]Host, error) {
	sqlStatement := `SELECT id, name, autoname, ip, snmp_community, status FROM host WHERE client_id=$1 order by autoname`
	var h Host
	var hosts []Host
	rows, err := c.Db.Query(sqlStatement, clientId)
	if err != nil {
		log.Warn(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&h.Id, &h.Name, &h.Autoname, &h.IP, &h.SNMP_community)
		if err != nil {
			err = fmt.Errorf("error fetchinghosts: %v", err)
			return nil, err
		}
		hosts = append(hosts, h)
	}

	err = rows.Err()
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	return hosts, nil
}

func (c configurator) PingParams() string {
	return "ping params"
}

func modelError(message string, err error) ModelError {
	if pgErr, ok := err.(*pq.Error); ok {
		return ModelError{OriginalMessage: pgErr.Message, OriginalHint: pgErr.Hint, Message: message}
	}
	return ModelError{OriginalMessage: err.Error(), Message: message}
}
