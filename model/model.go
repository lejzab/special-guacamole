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

func (pp PingProfile) String() string {
	return fmt.Sprintf("ID: %d, packet count: %d, packet interval: %d, packet size: %d, success: %d, timeout: %d.",
		pp.Id, pp.PacketCount, pp.PacketInterval, pp.PacketSize, pp.Success, pp.Timeout)
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
	sqlStatement := `select measure_type_id, name, default_value from measure_available_param
where measure_type_id in (select id from measure_type where family = 'PING')
and type = 'MEASUREMENT'
order by measure_type_id, name`

	var (
		profiles []PingProfile
		mt, dv   int
		n        string
		t        map[int]PingProfile
	)

	t = make(map[int]PingProfile)
	rows, err := c.Db.Query(sqlStatement)
	if err != nil {
		return nil, modelError("sql statemant error", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&mt, &n, &dv)
		p, ok := t[mt]
		if !ok {
			p = PingProfile{Id: mt}
		}
		switch n {
		case "packet_count":
			p.PacketCount = dv
		case "packet_interval":
			p.PacketInterval = dv
		case "packet_size":
			p.PacketSize = dv
		case "success":
			p.Success = dv
		case "timeout":
			p.Timeout = dv
		}
		t[mt] = p
		if err != nil {
			return nil, modelError("error fetching ping parameters", err)
		}
	}
	for key := range t {
		profiles = append(profiles, t[key])
	}
	log.Info(profiles)
	err = rows.Err()
	if err != nil {
		return nil, modelError("error fetching clients", err)
	}

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
