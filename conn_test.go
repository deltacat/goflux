package goflux

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const testDBName = "influx_test"

var client *Client

type InfluxTestSuite struct {
	suite.Suite
}

func (ts *InfluxTestSuite) SetupSuite() {
}

func (ts *InfluxTestSuite) SetupTest() {
	_ = client.DropDatabase(testDBName)
	_ = client.CreateDatabase(testDBName, true)
}

func (ts *InfluxTestSuite) TearDownTest() {
	//_ = client.DropDatabase(testDBName)
}

func TestViaSuite(t *testing.T) {
	conf := struct {
		Address, DatabaseName, RetentionPolicy, Precision string
	}{
		Address:      "http://localhost:8086", // todo: read conn and dbName from env
		DatabaseName: testDBName,
		RetentionPolicy: "",
		Precision:    Microsecond,
	}
	if c, err := CreateClient(conf.Address, "", "", conf.DatabaseName, conf.RetentionPolicy, conf.Precision); err != nil {
		panic(err)
	} else {
		client = c
	}
	defer client.Close()
	suite.Run(t, new(InfluxTestSuite))
}

func (ts *InfluxTestSuite) TestCreateDatabase() {
	ts.T().Run("Create Database", func(t *testing.T) {
		assert := require.New(t)
		tempDbName := "influx_test_temp"
		assert.NoError(client.CreateDatabase(tempDbName, false))
		assert.NoError(client.DropDatabase(tempDbName))
	})
}