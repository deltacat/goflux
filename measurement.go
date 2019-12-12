package goflux

import (
	"fmt"
)

// DropMeasurement drop a measurement with given name via query
func (c *Client) DropMeasurement(name string) error {
	cmd := fmt.Sprintf("DROP measurement %s", name)
	if _, err := c.queryEx(cmd); err != nil {
		return err
	}
	return nil
}

//ShowMeasurements show all measurements on the db
func (c *Client) ShowMeasurements(db string) []string {
	cmd := fmt.Sprintf("SHOW MEASUREMENTS ON %s", db)
	results, err := c.queryEx(cmd)
	if err != nil {
		return nil
	}
	if len(results.Series) == 0 || len(results.Series[0].Values) == 0 {
		return nil
	}
	values := results.Series[0].Values[0]
	measures := make([]string, len(values))
	for i, v := range values {
		measures[i] = v.(string)
	}
	return measures
}