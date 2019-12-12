package goflux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func (ts *InfluxTestSuite) TestDropMeasurement() {
	measureName := "test123"

	ts.T().Run("Create Measurement", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.WriteOne(measureName, Tags{"a":"a"}, Fields{"f": 1}, time.Now(), ""))
		measures := client.ShowMeasurements(client.db)
		assert.Equal(1, len(measures))
	})

	ts.T().Run("Drop Measurement", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.DropMeasurement(measureName))
		measures := client.ShowMeasurements(client.db)
		assert.Equal(0, len(measures))
	})
}
