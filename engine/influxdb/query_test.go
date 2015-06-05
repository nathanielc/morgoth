package influxdb_test

import (
	"fmt"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/engine/influxdb"
	"testing"
	"time"
)

func TestScheduleQueryShouldAddTimeWhereCondition(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		quertStr string
		start time.Time
		stop time.Time
		expectedStr string
	}

	testCases := []testCase{
		testCase{
			quertStr: "SELECT value FROM kbps",
			start: time.Unix(1433378783, 0),
			stop: time.Unix(1433379783, 0),
			expectedStr: "SELECT value FROM kbps WHERE time > \"2015-06-04 00:46:23\" AND time < \"2015-06-04 01:03:03\"",
		},
	}

	test := func(queryStr, start, stop, expectedStr) {
		q, err := influxdb.NewQuery(quertStr)
		assert.Nil(err)
		assert.Equal(
			expectedStr,
			q.QueryForTimeRange(start, stop),
		)
	}

}
