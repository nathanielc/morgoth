package kstest

import (
	log "github.com/cihub/seelog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"math"
	"sort"
	"time"
)

type fingerprint struct {
	Data []float64
	Count uint
}

type KSTest struct {
	reader       engine.Reader
	writer       engine.Writer
	config       *KSTestConf
	fingerprints []fingerprint
}

func (self *KSTest) Initialize(app app.App) error {
	self.reader = app.GetReader()
	self.writer = app.GetWriter()
	return nil
}

func (self *KSTest) Detect(metric metric.MetricID, start, stop time.Time) bool {
	log.Debugf("KSTest.Detect FP: %v", self.fingerprints)
	points := self.reader.GetData(metric, start, stop, -1)
	data := make([]float64, len(points))
	for i, point := range points {
		data[i] = point.Value
	}
	sort.Float64s(data)
	log.Debugf("Testing %v", data)

	minError := 0.0
	bestMatch := -1
	isMatch := false
	for i, fingerprint := range self.fingerprints {
		thresholdD := self.getThresholdD(len(fingerprint.Data), len(data))

		D := calcTestD(fingerprint.Data, data)
		log.Debug("D: ", D)
		if D < thresholdD {
			isMatch = true
		}
		e := (D - thresholdD) / thresholdD
		if bestMatch == -1 || e < minError {
			minError = e
			bestMatch = i
		}
	}

	anomalous := false
	if isMatch {
		anomalous = self.fingerprints[bestMatch].Count < self.config.NormalCount
		self.fingerprints[bestMatch].Count++
	} else {
		anomalous = true
		//We know its anomalous, now we need to update our fingerprints

		if len(self.fingerprints) == int(self.config.MaxFingerprints) {
			log.Debug("Reached MaxFingerprints")
			//TODO: Update bestMatch to learn new fingerprint
		} else {
			self.fingerprints = append(self.fingerprints, fingerprint{
				Data: data,
				Count: 1,
			})
		}
	}

	return anomalous
}
func (self *KSTest) getThresholdD(n, m int) float64 {
	c := 0.0
	switch self.config.Strictness {
	case 0: // 0.10
		c = 1.22
	case 1: // 0.05
		c = 1.36
	case 2: // 0.025
		c = 1.48
	case 3: // 0.01
		c = 1.63
	case 4: // 0.005
		c = 1.73
	case 5: // 0.001
		c = 1.95
	}
	return c * math.Sqrt(float64(n  + m) / float64( n * m))
}

func calcTestD(f1, f2 []float64) float64 {
	D := 0.0
	n := float64(len(f1))
	m := float64(len(f2))
	cdf1 := 0.0
	cdf2 := 0.0
	j := 0
	for _, x1 := range f1 {
		cdf1 += 1 / n
		for j < int(m) && x1 >= f2[j] {
			j++
			cdf2 += 1 / m
		}
		if d := math.Abs(cdf1 - cdf2); d > D {
			D = d
		}
		if j == int(m) { //Optimization only
			break
		}
	}
	return D
}
