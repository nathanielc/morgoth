package detectors_test

import (
	"testing"
	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/engine/generator"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/detector/mgof"
	"github.com/nvcook42/morgoth/schedule"
	"time"
	"math"
	"math/rand"
)


func init() {
	rand.Seed(42)
}


func m1(t int64) float64 {
	if t > 5 * periodSeconds && t < 5 * periodSeconds + 60 {
		return rand.Float64() * 0.2 - 0.1
	}
	return 2 * rand.Float64() * math.Sin(float64(t))
}

var rotation = schedule.Rotation{
	Period: 5*time.Minute,
	Resolution: time.Second,
}

var periodSeconds int64 = int64(rotation.Period/time.Second)

func TestMGOF1(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	factory := &mgof.MGOFFactory{}

	mgofConf := &mgof.MGOFConf{
		Min: -2,
		Max: 2,
		NullConfidence: 4,
	}
	functions := make(map[metric.MetricID]generator.Ft)

	functions["m1"] = m1

	detector, err := setup("m1", rotation, mgofConf, factory, functions)
	if !assert.Nil(err) {
		assert.Fail("Failed to create detector ", err.Error())
	}
	
	//Turn on mgof trace level logging so the results can be graphed
	//mgof.TraceLevel = 0

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
		float64(falsePositives) / float64(count) < 0.10,
		"Too many false positives %d",
		falsePositives,
	)

}
