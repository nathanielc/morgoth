package tukey_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	detector_test "github.com/nvcook42/morgoth/detector/test"
	"github.com/nvcook42/morgoth/detector/tukey"
	"github.com/nvcook42/morgoth/engine/generator"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"math"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(42)
}

func m1(t int64) float64 {
	if t > 5*periodSeconds && t < 5*periodSeconds+60 {
		return rand.Float64()*5.0 + 1
	}
	return 2 * rand.Float64() * math.Sin(float64(t))
}

var rotation = schedule.Rotation{
	Period:     5 * time.Minute,
	Resolution: time.Second,
}

var periodSeconds int64 = int64(rotation.Period / time.Second)

func TestTukey1(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	factory := &tukey.TukeyFactory{}

	tukeyConf := &tukey.TukeyConf{}
	tukeyConf.Default()
	functions := make(map[metric.MetricID]generator.Ft)

	functions["m1"] = m1

	detector, err := detector_test.SetupGeneratedDetectorTest("m1", rotation, tukeyConf, factory, functions)
	if !assert.Nil(err) {
		assert.Fail("Failed to create detector ", err.Error())
	}

	start := generator.TZero
	stop := start.Add(rotation.Period)
	count := 6
	falsePositives := 0
	for i := 0; i < count; i++ {
		anomalous := detector.Detect("m1", start, stop)
		if i == 5 {
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
