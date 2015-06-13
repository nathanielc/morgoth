package benchmarks

import (
	"github.com/nathanielc/morgoth/config"
	"testing"
)

func BenchmarkParseConfig(b *testing.B) {

	var data string = `---
engine:
  influxdb:
    host: 192.168.1.1
    port: 8086
    user: root
    password: root
    database: morgoth

schedule:
  rotations:
    - {period: 2m, resolution: 2s}
    - {period: 4m, resolution: 4s}
    - {period: 1d, resolution: 30m}
  delay: 15s

metrics:
  - pattern: cpu.*
    detectors:
      - mgof:
          min: 0
          max: 100
    notifiers:
  - pattern: .*
    detectors:
      - kstest: {}
    notifiers:


fittings:
  - rest:
      port: 7000
  - graphite: {}

logging:
    level: INFO
`
	for i := 0; i < b.N; i++ {
		config.Load([]byte(data))
	}
}
