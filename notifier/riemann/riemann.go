package riemann

import (
	"fmt"
	"github.com/bigdatadev/goryman"
	"github.com/golang/glog"
	metric "github.com/nvcook42/morgoth/metric/types"
	"time"
)

type RiemannNotifier struct {
	client *goryman.GorymanClient
}

func New(host string, port uint) *RiemannNotifier {

	rn := new(RiemannNotifier)

	rn.client = goryman.NewGorymanClient(fmt.Sprintf("%s:%d", host, port))
	err := rn.client.Connect()
	if err != nil {
		glog.Warningf("Error connecting to riemann: %s host: %s port: %d", err.Error(), host, port)
	}

	return rn
}

func (self *RiemannNotifier) Notify(detectorId string, metric metric.MetricID, start, stop time.Time) {

	err := self.client.SendEvent(&goryman.Event{
		Service:     string(metric),
		State:       "anomalous",
		Description: fmt.Sprintf("start: %v, stop: %v", start, stop),
		Tags:        []string{detectorId},
		Time:        start.Unix(),
		Metric:      int(stop.Sub(start)),
	})
	if err != nil {
		glog.Warning("Failed to send event ", err)
	}
}
