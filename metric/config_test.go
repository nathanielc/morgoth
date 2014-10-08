package metric_test

import (
	"github.com/nvcook42/morgoth/metric"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestMetricConfShouldParsePattern(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
pattern: .*
`

	mc := metric.MetricConf{}

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal(".*", mc.Pattern)

}
