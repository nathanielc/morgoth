package kstest

import (
	"fmt"
	"github.com/golang/glog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"github.com/nvcook42/morgoth/detector/metadata"
	"math"
	"encoding/json"
	"sort"
	"time"
)

type fingerprint struct {
	Data  []float64
	Count uint
}

type KSTest struct {
	rotation     schedule.Rotation
	reader       engine.Reader
	writer       engine.Writer
	config       *KSTestConf
	fingerprints map[metric.MetricID][]fingerprint
	meta         *metadata.MetadataStore
}

func (self *KSTest) GetIdentifier() string {
	return fmt.Sprintf(
		"%skstest_%d_%d_%d",
		self.rotation.GetPrefix(),
		self.config.Confidence,
		self.config.NormalCount,
		self.config.MaxFingerprints,
	)
}

func (self *KSTest) Initialize(app app.App, rotation schedule.Rotation) error {
	self.rotation = rotation
	self.reader = app.GetReader()
	self.writer = app.GetWriter()
	self.fingerprints = make(map[metric.MetricID][]fingerprint)

	meta, err := app.GetMetadataStore(self.GetIdentifier())
	if err != nil {
		return err
	}
	self.meta = meta

	return nil
}

func (self *KSTest) Detect(metric metric.MetricID, start, stop time.Time) bool {
	fingerprints, ok := self.fingerprints[metric]
	if !ok {
		fingerprints = self.load(metric)
	}

	points := self.reader.GetData(&self.rotation, metric, start, stop)
	data := make([]float64, len(points))
	for i, point := range points {
		data[i] = point.Value
	}
	sort.Float64s(data)

	minError := 0.0
	bestMatch := -1
	isMatch := false
	for i, fingerprint := range fingerprints {
		thresholdD := self.getThresholdD(len(fingerprint.Data), len(data))

		D := calcTestD(fingerprint.Data, data)
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
		anomalous = fingerprints[bestMatch].Count < self.config.NormalCount
		fingerprints[bestMatch].Count++
	} else {
		anomalous = true
		//We know its anomalous, now we need to update our fingerprints

		if len(fingerprints) == int(self.config.MaxFingerprints) {
			glog.V(2).Info("Reached MaxFingerprints")
			//TODO: Update bestMatch to learn new fingerprint
		} else {
			fingerprints = append(fingerprints, fingerprint{
				Data:  data,
				Count: 1,
			})
		}
	}

	self.fingerprints[metric] = fingerprints
	go self.save(metric)

	if glog.V(3) {
		jf, _ := json.Marshal(fingerprints)
		jd, _ := json.Marshal(data)
		glog.Infof(
			"%s|%s# { \"anomalous\" : %v, \"fingerprints\" : %s, \"current\" : %s }",
			self.GetIdentifier(),
			string(metric),
			anomalous,
			jf,
			jd,
		)
	}
	return anomalous
}

func (self *KSTest) getThresholdD(n, m int) float64 {
	c := 0.0
	switch self.config.Confidence {
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
	return c * math.Sqrt(float64(n+m)/float64(n*m))
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

func (self *KSTest) save(metric metric.MetricID) {

	data, err := json.Marshal(self.fingerprints[metric])
	if err != nil {
		glog.Error("Could not save KSTest", err.Error())
	}
	self.meta.StoreDoc(metric, data)
}

func (self *KSTest) load(metric metric.MetricID) []fingerprint {

	fingerprints := make([]fingerprint, 0)
	data := self.meta.GetDoc(metric)
	if len(data) != 0 {
		err := json.Unmarshal(data, &fingerprints)
		if err != nil {
			glog.Error("Could not load KSTest metadata", err.Error())
		}
	}
	return fingerprints
}
