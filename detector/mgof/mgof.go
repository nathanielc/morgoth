package mgof

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"github.com/nvcook42/morgoth/stat"
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
}

func (self *MGOF) Initialize(app app.App, rotation schedule.Rotation) error {
	self.rotation = rotation
	self.reader = app.GetReader()
	self.writer = app.GetWriter()

	self.threshold = stat.Xsquare_InvCDF(int64(self.config.NBins - 1))(self.config.CHI2)

	self.load()
	return nil
}

func (self *MGOF) Detect(metric metric.MetricID, start, stop time.Time) bool {
	fingerprints := self.fingerprints[metric]
	log.Debugf("MGOF.Detect Rotation: %s FP %v", self.rotation.GetPrefix(), fingerprints)
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
			log.Warnf("Not enough data points to trust histogram: %d < %d bins for metric %s", fingerprint.Hist.Count, nbins, string(metric))
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
			log.Debug("Reached MaxFingerprints")
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
	//go self.save()
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
	multiplier := 10.0
	fakeTotal := float64(hist.Count)*multiplier + float64(len(hist.Bins))
	empty := 1.0 / fakeTotal
	for i := range hist.Bins {
		hist.Bins[i] = empty + hist.Bins[i]*multiplier/fakeTotal
	}
}

func (self *MGOF) save(metric metric.MetricID) {

	data, err := json.Marshal(self.fingerprints[metric])
	if err != nil {
		log.Error("Could not save MGOF", err.Error())
	}
	self.writer.StoreDoc(self.rotation.GetPrefix()+"mgof."+string(metric), data)
}

func (self *MGOF) load() {

	data := self.reader.GetDoc(self.rotation.GetPrefix() + "mgof")
	if len(data) != 0 {
		err := json.Unmarshal(data, &self.fingerprints)
		if err != nil {
			log.Error("Could not load MGOF ", err.Error())
		}
	}
	if self.fingerprints == nil {
		self.fingerprints = make(map[metric.MetricID][]fingerprint)
	}
	for metric, fingerprints := range self.fingerprints {
		if len(fingerprints) > 0 &&
			len(fingerprints[0].Hist.Bins) != int(self.config.NBins) {
			self.fingerprints[metric] = make([]fingerprint, 0)
		}
	}
}
