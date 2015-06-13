package influxdb_test

import (
	"flag"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/engine/influxdb"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	defer glog.Flush()
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "1")
	}
	os.Exit(m.Run())
}

func TestScheduleQueryShouldAddTimeWhereCondition(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		quertStr    string
		start       time.Time
		stop        time.Time
		expectedStr string
	}

	testCases := []testCase{
		testCase{
			quertStr:    `SELECT value FROM kbps`,
			start:       time.Unix(1433378783, 0),
			stop:        time.Unix(1433379783, 0),
			expectedStr: `SELECT value FROM kbps WHERE time > '2015-06-04 00:46:23' AND time < '2015-06-04 01:03:03'`,
		},
		testCase{
			quertStr:    `SELECT value FROM kbps WHERE datacenter = nyc`,
			start:       time.Unix(1433378783, 0),
			stop:        time.Unix(1433379783, 0),
			expectedStr: `SELECT value FROM kbps WHERE datacenter = nyc AND time > '2015-06-04 00:46:23' AND time < '2015-06-04 01:03:03'`,
		},
	}

	test := func(tc testCase) {
		q, err := influxdb.NewQuery(tc.quertStr)
		assert.Nil(err)
		assert.Equal(
			tc.expectedStr,
			q.QueryForTimeRange(tc.start, tc.stop),
		)
	}

	for _, tc := range testCases {
		test(tc)
	}
}

func TestExecQuery(t *testing.T) {
	assert := assert.New(t)

	engine := new(influxdb.InfluxDBEngine)

	q, err := influxdb.NewQuery("SELECT value FROM /cpu_load_.*/ WHERE host='server01' group by *")
	assert.Nil(err)

	start := time.Unix(1433378783, 0)
	stop := time.Unix(1433624798, 0)
	qStr := q.QueryForTimeRange(start, stop)

	engine.ExecuteQuery(qStr)
}
