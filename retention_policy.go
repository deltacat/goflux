package goflux

import (
	"encoding/json"
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
	Name               string
	Duration           string
	Replication        int64
	Default            bool
	ShardGroupDuration string
}

var allPolicies = []RetentionPolicy{
	{Name: retentionOneWeek, Duration: durationOneWeek, Replication: 1},
	{Name: retentionOneMonth, Duration: durationOneMonth, Replication: 1},
	{Name: retentionThreeMonth, Duration: durationThreeMonth, Replication: 1},
	{Name: retentionOneYear, Duration: durationOneYear, Replication: 1},
}

// CreateAllRetentionPolicies create all application wanted retention policies
func (c *Client) CreateAllRetentionPolicies() error {
	return c.CreateRetentionPolicies(allPolicies)
}

//GetAllRetentionPolicies get all retention policies
func (c *Client) GetAllRetentionPolicies() ([]RetentionPolicy, error) {
	cmd := fmt.Sprintf("show retention policies")
	result, err := c.queryEx(cmd)
	if err != nil {
		return nil, err
	}
	series := result.Series
	if len(series) == 0 {
		return nil, ErrEmptyResults
	}

	columns := make(map[string]int)
	for i, col := range series[0].Columns {
		columns[col] = i
	}
	values := series[0].Values
	rps := make([]RetentionPolicy, len(values))
	for i, r := range values {
		replicaN, _ := r[columns["replicaN"]].(json.Number).Int64()
		rps[i] = RetentionPolicy{
			Name:               r[columns["name"]].(string),
			Duration:           r[columns["duration"]].(string),
			Replication:        replicaN,
			Default:            r[columns["default"]].(bool),
			ShardGroupDuration: r[columns["shardGroupDuration"]].(string),
		}
	}

	return rps, nil
}

//DropAllRetentionPolicies drop all retention policies which created by the app
//will not drop predefined policy "autogen"
//DANGER!!! this operation will permanently delete all measurements and data stored in the retention policy
func (c *Client) DropAllRetentionPolicies() error {
	return c.DropRetentionPolicies(allPolicies)
}

// CreateRetentionPolicy create a policy
func (c *Client) CreateRetentionPolicy(name string, duration string, replication int64) error {
	return c.CreateRetentionPolicies([]RetentionPolicy{{Name: name, Duration: duration, Replication: replication}})
}

// CreateRetentionPolicies create input retention policies
func (c *Client) CreateRetentionPolicies(rps []RetentionPolicy) error {
	if len(rps) == 0 {
		return nil
	}
	cmd := ""
	for _, rp := range rps {
		cmd += fmt.Sprintf("create retention policy %s on %s duration %s replication %d; ", rp.Name, c.db, rp.Duration, rp.Replication)
	}
	_, err := c.queryEx(cmd)
	return err
}

//SetDefaultRetentionPolicy change default retention policies
func (c *Client) SetDefaultRetentionPolicy(name string) error {
	if name == "" {
		return nil
	}
	cmd := fmt.Sprintf(`alter retention policy "%s" on %s default`, name, c.db)
	_, err := c.queryEx(cmd)
	return err
}

//GetDefaultRPName get default retention policy name
func (c *Client) GetDefaultRPName() string {
	results, err := c.GetAllRetentionPolicies()
	if err != nil {
		return ""
	}
	for _, r := range results {
		if r.Default {
			return r.Name
		}
	}
	return ""
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
		cmd += fmt.Sprintf("drop retention policy %s on %s;", rp.Name, c.db)
	}
	_, err := c.queryEx(cmd)
	return err
}

//DropRetentionPolicy drop a retention policy
//DANGER!!! this operation will permanently delete all measurements and data stored in the retention policy
func (c *Client) DropRetentionPolicy(name string) error {
	cmd := fmt.Sprintf("drop retention policy %s on %s", name, c.db)
	_, err := c.queryEx(cmd)
	return err
}
