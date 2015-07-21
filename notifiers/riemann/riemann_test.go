package riemann_test

import (
	"flag"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	metric "github.com/nathanielc/morgoth/metric/types"
	"github.com/nathanielc/morgoth/notifier/riemann"
	"testing"
	"time"
)

func init() {
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "1")
	}
}

func TestNotify(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	r, err := riemann.New("localhost", 5555)
	if !assert.Nil(err) {
		t.FailNow()
	}

	metric := metric.MetricID("cpu")
	stop := time.Now()
	start := stop.Add(-time.Hour)
	r.Notify("mgof", metric, start, stop)

	assert.True(false)
}
