package goflux

import (
	"fmt"
)

const (
	retentionPredefined = "autogen"
	retentionOneHour    = "a_hour"
	retentionOneDay     = "a_day"
	retentionOneWeek    = "a_week"
	retentionOneMonth   = "a_month"
	retentionThreeMonth = "tri_month"
	retentionOneYear    = "a_year"
	durationOneHour     = "1h"
	durationOneDay      = "1d"
	durationOneWeek     = "1w"
	durationOneMonth    = "30d"
	durationThreeMonth  = "90d"
	durationOneYear     = "365d"
)

// RetentionPolicy data struct for retention policy
type RetentionPolicy struct {
	name        string
	duration    string
	replication int
}

var allPolicies = []RetentionPolicy{
	{name: retentionOneWeek, duration: durationOneWeek, replication: 1},
	{name: retentionOneMonth, duration: durationOneMonth, replication: 1},
	{name: retentionThreeMonth, duration: durationThreeMonth, replication: 1},
	{name: retentionOneYear, duration: durationOneYear, replication: 1},
}

// CreateAllRetentionPolicies create all application wanted retention policies
func (c *Client) CreateAllRetentionPolicies() error {
	return c.CreateRetentionPolicies(allPolicies)
}

//GetAllRetentionPolicies get all retention policies
func (c *Client) GetAllRetentionPolicies() ([][]interface{}, error) {
	cmd := fmt.Sprintf("show retention policies")
	result, err := c.query(cmd)
	if err != nil {
		return nil, err
	}
	series := result.Series
	if len(series) == 0 {
		return nil, ErrEmptyResults
	}
	return series[0].Values, nil
}

//DropAllRetentionPolicies drop all retention policies which created by the app
//will not drop predefined policy "autogen"
//DANGER!!! this operation will permanently delete all measurements and data stored in the retention policy
func (c *Client) DropAllRetentionPolicies() error {
	return c.DropRetentionPolicies(allPolicies)
}

// CreateRetentionPolicy create a policy
func (c *Client) CreateRetentionPolicy(name string, duration string, replication int) error {
	return c.CreateRetentionPolicies([]RetentionPolicy{{name: name, duration: duration, replication: replication}})
}

// CreateRetentionPolicies create input retention policies
func (c *Client) CreateRetentionPolicies(rps []RetentionPolicy) error {
	if len(rps) == 0 {
		return nil
	}
	cmd := ""
	for _, rp := range rps {
		cmd += fmt.Sprintf("create retention policy %s on %s duration %s replication %d; ", rp.name, c.db, rp.duration, rp.replication)
	}
	_, err := c.query(cmd)
	return err
}

//DropRetentionPolicies drop all retention policies which created by the app
//will not drop predefined policy "autogen"
//DANGER!!! this operation will permanently delete all measurements and data stored in the retention policy
func (c *Client) DropRetentionPolicies(rps []RetentionPolicy) error {
	if len(rps) == 0 {
		return nil
	}
	cmd := ""
	for _, rp := range rps {
		cmd += fmt.Sprintf("drop retention policy %s on %s;", rp.name, c.db)
	}
	_, err := c.query(cmd)
	return err
}

//DropRetentionPolicy drop a retention policy
//DANGER!!! this operation will permanently delete all measurements and data stored in the retention policy
func (c *Client) DropRetentionPolicy(name string) error {
	cmd := fmt.Sprintf("drop retention policy %s on %s", name, c.db)
	_, err := c.query(cmd)
	return err
}
