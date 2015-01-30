package riemann

import (
	"fmt"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/amir/raidman"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	metric "github.com/nvcook42/morgoth/metric/types"
	"time"
)

const (
	maxRetries = 5
)

type RiemannNotifier struct {
	client  *raidman.Client
	host    string
	port    uint
	retries int
}

func New(host string, port uint) (*RiemannNotifier, error) {

	rn := &RiemannNotifier{
		host: host,
		port: port,
	}
	err := rn.connect()
	if err != nil {
		return nil, err
	}
	return rn, nil
}

func (self *RiemannNotifier) connect() error {
	if self.client != nil {
		self.client.Close()
	}
	client, err := raidman.Dial("tcp", fmt.Sprintf("%s:%d", self.host, self.port))
	if err != nil {
		glog.Warningf(
			"Error connecting to riemann: %s host: %s port: %d",
			err.Error(),
			self.host,
			self.port,
		)
		return err
	}
	self.client = client
	self.retries = maxRetries
	return nil
}

func (self *RiemannNotifier) Notify(detectorId string, metric metric.MetricID, start, stop time.Time) {
	event := &raidman.Event{
		Service:     string(metric),
		Host:        "morgoth",
		State:       "anomalous",
		Description: fmt.Sprintf("start: %v, stop: %v", start, stop),
		Tags:        []string{detectorId},
		Time:        start.Unix(),
		Metric:      int(stop.Sub(start)),
		Ttl:         30,
	}
	err := self.client.Send(event)
	if err != nil {
		err := self.connect()
		if err == nil && self.retries > 0 {
			self.retries--
			glog.Warning("Failed to send event retrying...")
			self.Notify(detectorId, metric, start, stop)
		} else {
			glog.Error("Failed to send event ", err)
		}
	}
}
