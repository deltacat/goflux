package goflux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type simpleData struct {
	City      string
	Brand     string
	Customers int
	Sales     int
	Time      time.Time
}

const simpleMeasurement = "simple"

func addOneSalesData(data simpleData) error {
	tags := Tags{
		"city":  data.City,
		"brand": data.Brand,
	}
	fields := Fields{
		"customers": data.Customers,
		"sales":     data.Sales,
	}
	return client.WriteOne(simpleMeasurement, tags, fields, data.Time)
}

func fetchRecentSalesData(duration time.Duration) (*Result, error) {
	return client.FetchRecent(simpleMeasurement, duration)
}

func buildTestData() []simpleData {
	var dataArr []simpleData
	timestamp := time.Now().Add(-time.Second)
	for i := 0; i < 10; i++ {
		dataArr = append(dataArr, simpleData{
			City:      "Hangzhou",
			Brand:     "deltacat",
			Customers: 2,
			Sales:     3,
			Time:      timestamp,
		})
		timestamp = timestamp.Add(time.Nanosecond)
	}
	return dataArr
}

func (ts *InfluxTestSuite) TestInsertOneData() {
	tests := buildTestData()

	ts.T().Run("insert data", func(t *testing.T) {
		assert := require.New(t)
		for _, data := range tests {
			assert.NoError(addOneSalesData(data))
		}
	})

	ts.T().Run("fetch recent", func(t *testing.T) {
		assert := require.New(t)
		result, err := fetchRecentSalesData(time.Minute)
		assert.NoError(err)
		assert.Equal(1, len(result.Series))
		assert.Equal(10, len(result.Series[0].Values))
	})
}

func (ts *InfluxTestSuite) TestInsertMultiPoints() {
	tests := buildTestData()

	ts.T().Run("insert multi points", func(t *testing.T) {
		assert := require.New(t)
		points := make([]*Point, 0)
		for _, data := range tests {
			tags := Tags{
				"city":  data.City,
				"brand": data.Brand,
			}
			fields := Fields{
				"customers": data.Customers,
				"sales":     data.Sales,
			}
			pt, err := NewPoint(simpleMeasurement, tags, fields, data.Time)
			assert.NoError(err)
			points = append(points, pt)
		}
		assert.Equal(10, len(points))
		assert.NoError(client.WriteAll(points, "", ""))
	})

	ts.T().Run("fetch recent", func(t *testing.T) {
		assert := require.New(t)
		result, err := fetchRecentSalesData(time.Minute)
		assert.NoError(err)
		assert.Equal(1, len(result.Series))
		assert.Equal(10, len(result.Series[0].Values))
	})
}
