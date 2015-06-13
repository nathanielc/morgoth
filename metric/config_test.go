package metric_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/metric"
	"testing"
)

func TestMetricSupervisorConfShouldParsePattern(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
pattern: .*
`

	mc := metric.MetricSupervisorConf{}

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal(".*", mc.Pattern)

}
