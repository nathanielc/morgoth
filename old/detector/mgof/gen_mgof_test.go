package mgof_test

import (
	"flag"
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/detector"
	"github.com/nathanielc/morgoth/detector/mgof"
	detector_test "github.com/nathanielc/morgoth/detector/test"
	"github.com/nathanielc/morgoth/engine/generator"
	metric "github.com/nathanielc/morgoth/metric/types"
	"github.com/nathanielc/morgoth/schedule"
	"math"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(42)
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "1")
		flag.Set("v", fmt.Sprintf("%d", detector.TraceLevel))
	}
}

func m1(t int64) float64 {
	if t > 5*periodSeconds && t < 5*periodSeconds+60 {
		return rand.Float64()*0.2 - 0.1
	}
	return 2 * rand.Float64() * math.Sin(float64(t))
}

var rotation = schedule.Rotation{
	Period:     5 * time.Minute,
	Resolution: time.Second,
}

var periodSeconds int64 = int64(rotation.Period / time.Second)

func TestMGOF1(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	factory := &mgof.MGOFFactory{}

	mgofConf := &mgof.MGOFConf{
		Min:            -2,
		Max:            2,
		NullConfidence: 4,
	}
	functions := make(map[metric.MetricID]generator.Ft)

	functions["m1"] = m1

	detector, err := detector_test.SetupGeneratedDetectorTest("m1", rotation, mgofConf, factory, functions)
	if !assert.Nil(err) {
		assert.Fail("Failed to create detector ", err.Error())
	}

	start := generator.TZero
	stop := start.Add(rotation.Period)
	count := 50
	falsePositives := 0
	for i := 0; i < count; i++ {
		anomalous := detector.Detect("m1", start, stop)
		if i < 3 || i == 5 {
			assert.True(anomalous, "Run %d", i)
		} else {
			if anomalous {
				falsePositives++
			}
		}
		start = stop
		stop = stop.Add(rotation.Period)
	}

	assert.True(
		float64(falsePositives)/float64(count) < 0.10,
		"Too many false positives %d",
		falsePositives,
	)
}
