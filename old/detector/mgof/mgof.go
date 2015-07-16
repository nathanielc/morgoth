package mgof

import (
	"encoding/json"
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	app "github.com/nathanielc/morgoth/app/types"
	"github.com/nathanielc/morgoth/detector"
	"github.com/nathanielc/morgoth/detector/metadata"
	"github.com/nathanielc/morgoth/engine"
	metric "github.com/nathanielc/morgoth/metric/types"
	"github.com/nathanielc/morgoth/schedule"
	"github.com/nathanielc/morgoth/stat"
	"math"
	"time"
)

type fingerprint struct {
	Hist  *engine.Histogram
	Count uint
}

type MGOF struct {
	rotation     schedule.Rotation
	reader       engine.Reader
	writer       engine.Writer
	config       *MGOFConf
	fingerprints map[metric.MetricID][]fingerprint
	threshold    float64
	meta         metadata.MetadataStore
}

func (self *MGOF) GetIdentifier() string {
	return fmt.Sprintf(
		"%smgof_%0.4f_%0.4f_%d_%d_%d_%d",
		self.rotation.GetPrefix(),
		self.config.Min,
		self.config.Max,
		self.config.NullConfidence,
		self.config.NBins,
		self.config.NormalCount,
		self.config.MaxFingerprints,
	)
}

func (self *MGOF) Initialize(app app.App, rotation schedule.Rotation) error {
	self.rotation = rotation
	self.reader = app.GetReader()
	self.writer = app.GetWriter()
	self.fingerprints = make(map[metric.MetricID][]fingerprint)

	// n the null confidence is the number of 9s in chi2
	n := int(self.config.NullConfidence)
	chi2 := (math.Pow10(n) - 1.0) / math.Pow10(n)
	self.threshold = stat.Xsquare_InvCDF(int64(self.config.NBins - 1))(chi2)

	meta, err := app.GetMetadataStore(self.GetIdentifier())
	if err != nil {
		return err
	}
	self.meta = meta

	return nil
}

func (self *MGOF) Detect(metric metric.MetricID, start, stop time.Time) bool {
	fingerprints, ok := self.fingerprints[metric]
	if !ok {
		fingerprints = self.load(metric)
	}
	nbins := self.config.NBins
	hist := self.reader.GetHistogram(
		&self.rotation,
		metric,
		nbins,
		start,
		stop,
		self.config.Min,
		self.config.Max,
	)

	fillEmptyValues(hist)

	minRE := 0.0
	bestMatch := -1
	isMatch := false
	for i, fingerprint := range fingerprints {
		if fingerprint.Hist.Count < nbins {
			glog.Warningf("Not enough data points to trust histogram: %d < %d bins for metric %s", fingerprint.Hist.Count, nbins, string(metric))
			continue
		}

		re := relativeEntropy(hist, fingerprint.Hist)
		if float64(2*hist.Count)*re < self.threshold {
			isMatch = true
		}
		if bestMatch == -1 || re < minRE {
			minRE = re
			bestMatch = i
		}
	}

	anomalous := false
	if isMatch {
		anomalous = fingerprints[bestMatch].Count < self.config.NormalCount
		fingerprints[bestMatch].Count++
	} else {
		anomalous = true
		//We know whether its anomalous, now we need to update our fingerprints

		if len(fingerprints) == int(self.config.MaxFingerprints) {
			glog.V(2).Info("Reached MaxFingerprints")
			//TODO: Update bestMatch to learn new fingerprint
			ratio := 1 / float64(fingerprints[bestMatch].Count)
			for i, p := range fingerprints[bestMatch].Hist.Bins {
				fingerprints[bestMatch].Hist.Bins[i] = (1-ratio)*p + ratio*hist.Bins[i]
			}
		} else {
			fingerprints = append(fingerprints, fingerprint{
				Hist:  hist,
				Count: 1,
			})
		}
	}

	self.fingerprints[metric] = fingerprints
	go self.save(metric)

	if glog.V(detector.TraceLevel) {
		jf, _ := json.Marshal(fingerprints)
		jh, _ := json.Marshal(hist)
		glog.Infof(
			"%s|%s# { \"anomalous\" : %v, \"fingerprints\" : %s, \"current\" : %s }",
			self.GetIdentifier(),
			string(metric),
			anomalous,
			jf,
			jh,
		)
	}
	return anomalous
}

func relativeEntropy(q, p *engine.Histogram) float64 {
	entropy := 0.0
	for i := range q.Bins {
		entropy += q.Bins[i] * math.Log(q.Bins[i]/p.Bins[i])
	}
	return entropy
}

func fillEmptyValues(hist *engine.Histogram) {
	multiplier := 100.0
	count := float64(hist.Count)
	fakeTotal := count*multiplier + float64(len(hist.Bins))
	empty := 1.0 / fakeTotal
	for i := range hist.Bins {
		hist.Bins[i] = empty + hist.Bins[i]*count*multiplier/fakeTotal
	}
}

func (self *MGOF) save(metric metric.MetricID) {

	data, err := json.Marshal(self.fingerprints[metric])
	if err != nil {
		glog.Error("Could not save MGOF metadata", err.Error())
	}
	self.meta.StoreDoc(metric, data)
}

func (self *MGOF) load(metric metric.MetricID) []fingerprint {

	fingerprints := make([]fingerprint, 0)
	data := self.meta.GetDoc(metric)
	if len(data) != 0 {
		err := json.Unmarshal(data, &fingerprints)
		if err != nil {
			glog.Error("Could not load MGOF metadata ", err.Error())
		}
	}

	return fingerprints
}
