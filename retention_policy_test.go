package goflux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func checkAddedPolicies(assert *require.Assertions, total int) string {
	results, err := client.GetAllRetentionPolicies()
	assert.NoError(err)
	assert.Equal(total, len(results))
	return results[total-1][0].(string)
}

var tests = []RetentionPolicy{
	{name: retentionOneHour, duration: durationOneHour, replication: 1},
	{name: retentionOneDay, duration: durationOneDay, replication: 1},
	{name: retentionOneWeek, duration: durationOneWeek, replication: 1},
	{name: retentionOneMonth, duration: durationOneMonth, replication: 1},
	{name: retentionOneYear, duration: durationOneYear, replication: 1},
}

func (ts *InfluxTestSuite) TestCreateRetentionPolicy() {

	ts.T().Run("read default", func(t *testing.T) {
		assert := require.New(t)
		assert.Equal(retentionPredefined, checkAddedPolicies(assert, 1))
	})

	ts.T().Run("add policies", func(t *testing.T) {
		assert := require.New(t)
		for i, tt := range tests {
			assert.NoError(client.CreateRetentionPolicy(tt.name, tt.duration, tt.replication))
			assert.Equal(tt.name, checkAddedPolicies(assert, i+2))
		}
	})

	ts.T().Run("add exists", func(t *testing.T) {
		assert := require.New(t)
		for _, tt := range tests {
			assert.NoError(client.CreateRetentionPolicy(tt.name, tt.duration, tt.replication))
			assert.Equal(tests[len(tests)-1].name, checkAddedPolicies(assert, len(tests)+1))
		}
	})

	ts.T().Run("add exists diff parameter", func(t *testing.T) {
		assert := require.New(t)
		for _, tt := range tests {
			assert.Error(client.CreateRetentionPolicy(tt.name, "13d", tt.replication))
		}
		results, err := client.GetAllRetentionPolicies()
		assert.NoError(err)
		assert.Equal(len(tests)+1, len(results))
		assert.Equal(retentionPredefined, results[0][0])
		for i, tt := range tests {
			assert.Equal(tt.name, results[i+1][0])
		}
	})

}

func (ts *InfluxTestSuite) TestCreateAllRetentionPolicies() {
	ts.T().Run("create all", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.CreateAllRetentionPolicies())
		addedNumber := len(allPolicies)
		//firstName := retentionPredefined
		lastName := allPolicies[addedNumber-1].name
		assert.Equal(lastName, checkAddedPolicies(assert, addedNumber+1))
	})
}

func (ts *InfluxTestSuite) TestDropRetentionPolicies() {
	ts.T().Run("drop one", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.CreateRetentionPolicies(tests))
		rpNumbers := len(tests)
		for i, tt := range tests {
			assert.NoError(client.DropRetentionPolicy(tt.name))
			results, err := client.GetAllRetentionPolicies()
			assert.NoError(err)
			assert.Equal(rpNumbers-i, len(results))
		}
	})
	ts.T().Run("drop all", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.CreateAllRetentionPolicies())
		assert.NoError(client.DropAllRetentionPolicies())
		results, err := client.GetAllRetentionPolicies()
		assert.NoError(err)
		assert.Equal(1, len(results))
		assert.Equal(retentionPredefined, results[0][0])
	})
}
