package riemann

import (
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/amir/raidman"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
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

func (self *RiemannNotifier) Notify(msg string, w *morgoth.Window) {
	tags := make([]string, 0, len(w.Tags))
	for t, v := range w.Tags {
		tags = append(tags, fmt.Sprintf("%s=%s", t, v))
	}
	event := &raidman.Event{
		Service:     w.Name,
		Host:        "morgoth",
		State:       "anomalous",
		Description: fmt.Sprintf("%s start: %v, stop: %v", msg, w.Start, w.Stop),
		Tags:        tags,
		Time:        w.Start.Unix(),
		Metric:      int(w.Stop.Sub(w.Start)),
		Ttl:         30,
	}
	err := self.client.Send(event)
	if err != nil {
		glog.Error("Failed to send event Riemann", err)
	}
}
