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
func (c *Client) WriteOne(name string, tags Tags, fields Fields, ts time.Time, rp string) error {
	// create a point
	pt, err := influx.NewPoint(name, tags, fields, ts)
	if err != nil {
		log.WithError(err).Error("create point")
		return err
	}

	return c.WriteAll([]*Point{pt}, rp)
}

// WriteAll write all points at once.
func (c *Client) WriteAll(pts []*Point, rp string) error {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:        c.db,
		Precision:       c.pr,
		RetentionPolicy: c.rp,
	})
	if err != nil {
		log.WithError(err).Error("new batch point")
		return err
	}
	if rp != "" {
		bp.SetRetentionPolicy(rp)
	}

	bp.AddPoints(pts)

	// Write the batch
	return c.Write(bp)
}

// FetchRecent fetch recent data in <duration> from measurement <name>
func (c *Client) FetchRecent(name string, duration time.Duration) (*Result, error) {
	cmd := fmt.Sprintf("select * from %s where time > now() - %dns", name, duration.Nanoseconds())
	return c.queryEx(cmd)
}

//GetLastPoint get latest point of a measurement
func (c *Client) GetLastPoint(name string, tags Tags) (*Result, error) {
	whereClause := ""
	for k, v := range tags {
		if whereClause == "" {
			whereClause = fmt.Sprintf("where %s = '%s'", k, v)
		} else {
			whereClause = fmt.Sprintf("%s and %s = '%s'", whereClause, k, v)
		}
	}
	cmd := fmt.Sprintf("select * from %s %s order by time desc limit 1", name, whereClause)
	return c.queryEx(cmd)
}
