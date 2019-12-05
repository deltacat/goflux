package goflux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func checkAddedPolicies(assert *require.Assertions, total int) string {
	results, err := client.GetAllRetentionPolicies()
	assert.NoError(err)
	assert.Equal(total, len(results))
	return results[total-1].Name
}

var tests = []RetentionPolicy{
	{Name: retentionOneHour, Duration: durationOneHour, Replication: 1},
	{Name: retentionOneDay, Duration: durationOneDay, Replication: 1},
	{Name: retentionOneWeek, Duration: durationOneWeek, Replication: 1},
	{Name: retentionOneMonth, Duration: durationOneMonth, Replication: 1},
	{Name: retentionOneYear, Duration: durationOneYear, Replication: 1},
}

func (ts *InfluxTestSuite) TestCreateRetentionPolicy() {

	ts.T().Run("read default", func(t *testing.T) {
		assert := require.New(t)
		assert.Equal(retentionPredefined, checkAddedPolicies(assert, 1))
	})

	ts.T().Run("add policies", func(t *testing.T) {
		assert := require.New(t)
		for i, tt := range tests {
			assert.NoError(client.CreateRetentionPolicy(tt.Name, tt.Duration, tt.Replication))
			assert.Equal(tt.Name, checkAddedPolicies(assert, i+2))
		}
	})

	ts.T().Run("add exists", func(t *testing.T) {
		assert := require.New(t)
		for _, tt := range tests {
			assert.NoError(client.CreateRetentionPolicy(tt.Name, tt.Duration, tt.Replication))
			assert.Equal(tests[len(tests)-1].Name, checkAddedPolicies(assert, len(tests)+1))
		}
	})

	ts.T().Run("add exists diff parameter", func(t *testing.T) {
		assert := require.New(t)
		for _, tt := range tests {
			assert.Error(client.CreateRetentionPolicy(tt.Name, "13d", tt.Replication))
		}
		results, err := client.GetAllRetentionPolicies()
		assert.NoError(err)
		assert.Equal(len(tests)+1, len(results))
		assert.Equal(retentionPredefined, results[0].Name)
		for i, tt := range tests {
			assert.Equal(tt.Name, results[i+1].Name)
		}
	})

}

func (ts *InfluxTestSuite) TestCreateAllRetentionPolicies() {
	ts.T().Run("create all", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.CreateAllRetentionPolicies())
		addedNumber := len(allPolicies)
		//firstName := retentionPredefined
		lastName := allPolicies[addedNumber-1].Name
		assert.Equal(lastName, checkAddedPolicies(assert, addedNumber+1))
	})
}

func (ts *InfluxTestSuite) TestDropRetentionPolicies() {
	ts.T().Run("drop one", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.CreateRetentionPolicies(tests))
		rpNumbers := len(tests)
		for i, tt := range tests {
			assert.NoError(client.DropRetentionPolicy(tt.Name))
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
		assert.Equal(retentionPredefined, results[0].Name)
	})
}

func (ts *InfluxTestSuite) TestSetDefaultRetentionPolicy() {
	ts.T().Run("set default rp", func(t *testing.T) {
		assert := require.New(t)
		assert.NoError(client.CreateAllRetentionPolicies())

		rp := ""
		assert.NoError(client.SetDefaultRetentionPolicy(rp))
		assert.Equal(retentionPredefined, client.GetDefaultRPName())

		rp = retentionOneWeek
		assert.NoError(client.SetDefaultRetentionPolicy(rp))
		assert.Equal(rp, client.GetDefaultRPName())

		rp = retentionOneYear
		assert.NoError(client.SetDefaultRetentionPolicy(rp))
		assert.Equal(rp, client.GetDefaultRPName())

		rp = ""
		assert.NoError(client.SetDefaultRetentionPolicy(rp))
		assert.Equal(retentionOneYear, client.GetDefaultRPName())

		rp = retentionPredefined
		assert.NoError(client.SetDefaultRetentionPolicy(rp))
		assert.Equal(rp, client.GetDefaultRPName())

	})
}
