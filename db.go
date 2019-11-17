package influx

import (
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influx "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

type (
	// Client influx client instance and necessary parameters
	Client struct {
		cli       influx.Client
		db        string
		precision string
	}
)

//CreateClient create an influx client holder
func CreateClient(addr, user, pass, db, precision string) (*Client, error) {
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

// NewPoint create a new point. this is wrapping influx.NewPoint()
func NewPoint(name string, tags Tags, fields Fields, ts time.Time) (*Point, error) {
	// create a point
	return influx.NewPoint(name, tags, fields, ts)
}

// WriteOne write one point to measurement <name>
func (c *Client) WriteOne(name string, tags Tags, fields Fields, ts time.Time) error {
	// create a point
	pt, err := influx.NewPoint(name, tags, fields, ts)
	if err != nil {
		log.WithError(err).Error("create point")
		return err
	}

	return c.WriteAll([]*Point{pt}, "", "")
}

// WriteAll write all points at once.
func (c *Client) WriteAll(pts []*Point, rp string, tp string) error {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  c.db,
		Precision: c.precision,
	})
	if err != nil {
		log.WithError(err).Error("new batch point")
		return err
	}
	if rp != "" {
		bp.SetRetentionPolicy(rp)
	}
	if tp != "" {
		_ = bp.SetPrecision(tp)
	}

	bp.AddPoints(pts)

	// Write the batch
	return c.cli.Write(bp)
}

// FetchRecent fetch recent data in <duration> from measurement <name>
func (c *Client) FetchRecent(name string, duration time.Duration) (*Result, error) {
	//cmd := fmt.Sprintf("select * from %s where time > now() - %dns", name, duration)
	cmd := fmt.Sprintf("select * from %s", name)
	return c.query(cmd)
}

func (c *Client) GetLastPoint(name string, tags Tags) (*Result, error) {
	tagsC := ""
	for k, v := range tags {
		c := ""
		if tagsC != "" {
			c = " and"
		}
		tagsC += fmt.Sprintf("%s %s = '%s'", c, k, v)
	}
	cmd := fmt.Sprintf("select last(*) from %s where %s", name, tagsC)
	return c.query(cmd)
}

//Close close the influx instance
func (c *Client) Close() {
	if c.cli != nil {
		if err := c.cli.Close(); err != nil {
			log.WithError(err).Error("close connection")
		}
	}
}

// single command query
func (c *Client) query(command string) (*Result, error) {
	q := influx.NewQuery(command, c.db, c.precision)
	response, err := c.cli.Query(q)
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
