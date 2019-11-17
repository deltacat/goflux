package goflux

import (
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

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

//GetLastPoint get latest point of a measurement
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
