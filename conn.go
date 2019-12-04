package goflux

import (
	"fmt"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influx "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

type (
	// Client influx client instance and necessary parameters
	Client struct {
		cli       influx.Client
		db        string
		rp        string
		precision string
	}
)

//CreateClient create an influx client holder
func CreateClient(addr, user, pass, db, rp, precision string) (*Client, error) {
	cli, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return nil, err
	}
	c := &Client{
		cli:       cli,
		db:        db,
		rp:        rp,
		precision: precision,
	}
	if err := c.CreateDatabase(db, true); err != nil {
		return nil, err
	}
	return c, c.CreateAllRetentionPolicies()
}

// CreateDatabase create a Database with a query
func (c *Client) CreateDatabase(name string, use bool) error {
	cmd := fmt.Sprintf("CREATE DATABASE %s", name)
	if _, err := c.query(cmd); err != nil {
		return err
	}
	if use {
		c.db = name
	}
	return nil
}

// UseDatabase set database with input name as current using database
func (c *Client) UseDatabase(name string) {
	c.db = name
}

// DropDatabase drop a Database with via query
func (c *Client) DropDatabase(name string) error {
	cmd := fmt.Sprintf("DROP DATABASE %s", name)
	if _, err := c.query(cmd); err != nil {
		return err
	}
	return nil
}

//Close close the influx instance
func (c *Client) Close() {
	if c.cli != nil {
		if err := c.cli.Close(); err != nil {
			log.WithError(err).Error("close connection")
		}
	}
}
