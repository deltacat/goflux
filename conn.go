package goflux

import (
	"fmt"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Client influx client instance and necessary parameters
type Client struct {
	influx.Client
	db, rp, pr string // database, retentionPolicy, precision
}

//CreateClient create an influx client holder
func CreateClient(addr, user, pass, db, rp, pr string) (*Client, error) {
	cli, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return nil, err
	}
	c := &Client{
		Client: cli,
		db:     db,
		rp:     rp,
		pr:     pr,
	}
	if err := c.CreateDatabase(db, true); err != nil {
		return nil, err
	}
	return c, c.CreateAllRetentionPolicies()
}

// CreateDatabase create a Database with a query
func (c *Client) CreateDatabase(name string, use bool) error {
	cmd := fmt.Sprintf("CREATE DATABASE %s", name)
	if _, err := c.queryEx(cmd); err != nil {
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
	if _, err := c.queryEx(cmd); err != nil {
		return err
	}
	return nil
}

// single command query
func (c *Client) queryEx(command string) (*Result, error) {
	q := influx.NewQuery(command, c.db, c.pr)
	response, err := c.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	results := response.Results
	if len(results) == 0 {
		return nil, ErrEmptyResults
	}
	return &results[0], nil
}
