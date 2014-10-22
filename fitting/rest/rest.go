package rest

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"net"
	"net/http"
	"strconv"
	"time"
)

type RESTFitting struct {
	port     uint
	reader   engine.Reader
	writer   engine.Writer
	listener net.Listener
	handler  rest.ResourceHandler
}

func (self *RESTFitting) Name() string {
	return "REST"
}

func (self *RESTFitting) Start(app app.App) {
	self.reader = app.GetReader()
	self.writer = app.GetWriter()
	self.handler = rest.ResourceHandler{
		EnableStatusService: true,
	}
	self.handler.SetRoutes(
		&rest.Route{"GET", "/status", self.status},
		&rest.Route{"GET", "/metrics", self.metrics},
		&rest.Route{"GET", "/data/:metric", self.metricData},
		&rest.Route{"GET", "/anomalies/:metric", self.anomalies},
		//&rest.Route{"POST", "/detect/:metric", self.detect},
		&rest.Route{"DELETE", "/delete/:metric", self.delete},
	)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", self.port))
	if err != nil {
		log.Error("Error starting REST fitting %s", err.Error())
		return
	}
	self.listener = listener
	err = http.Serve(self.listener, &self.handler)
	if err != nil {
		log.Debug(err)
		return
	}
	log.Info("REST fitting is done")
}

func (self *RESTFitting) Stop() {
	self.listener.Close()
}

////////////////////////////////
// Parameter Parsing/Validation
////////////////////////////////

func getTime(req *rest.Request, name string) (time.Time, error) {

	value := req.URL.Query().Get(name)
	if len(value) == 0 {
		return time.Time{}, nil
	}

	timestamp, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(timestamp, 0).UTC()

	return tm, nil
}

////////////////////////////////
// Hanlders
////////////////////////////////

func (self *RESTFitting) status(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(self.handler.GetStatus())
}

func (self *RESTFitting) metrics(w rest.ResponseWriter, req *rest.Request) {
	metrics := self.reader.GetMetrics()
	data := make(map[string]interface{}, 1)
	data["metrics"] = metrics
	w.WriteJson(data)
}

func (self *RESTFitting) metricData(w rest.ResponseWriter, req *rest.Request) {

	metric := metric.MetricID(req.PathParam("metric"))
	start, err := getTime(req, "start")
	if err != nil {
		rest.Error(w, "Could not parse 'start'"+err.Error(), 400)
		return
	}
	log.Debugf("Start: %v", start)

	stop, err := getTime(req, "stop")
	if err != nil {
		rest.Error(w, "Could not parse 'stop'"+err.Error(), 400)
		return
	}
	log.Debugf("Stop: %v", stop)

	data := make(map[string]interface{}, 1)
	data["metric"] = metric
	points := self.reader.GetData(metric, start, stop, 0)
	formatedPoints := make([][]interface{}, len(points))
	for i, point := range points {
		formatedPoints[i] = []interface{}{point.Time.Format(time.RFC3339Nano), point.Value}
	}
	data["data"] = formatedPoints

	w.WriteJson(data)
}

func (self *RESTFitting) anomalies(w rest.ResponseWriter, req *rest.Request) {

	metric := metric.MetricID(req.PathParam("metric"))

	start, err := getTime(req, "start")
	if err != nil {
		rest.Error(w, "Could not parse 'start'"+err.Error(), 400)
		return
	}
	log.Debugf("Start: %v", start)

	stop, err := getTime(req, "stop")
	if err != nil {
		rest.Error(w, "Could not parse 'stop'"+err.Error(), 400)
		return
	}
	log.Debugf("Stop: %v", stop)

	data := make(map[string]interface{}, 1)
	data["metric"] = metric
	anomalies := self.reader.GetAnomalies(metric, start, stop)
	formatedAnomalies := make([]map[string]interface{}, len(anomalies))
	for i, anomaly := range anomalies {
		formatedAnomalies[i] = map[string]interface{}{
			"id":    anomaly.UUID.String(),
			"start": anomaly.Start.Format(time.RFC3339Nano),
			"stop":  anomaly.Stop.Format(time.RFC3339Nano),
		}
	}
	data["data"] = formatedAnomalies

	w.WriteJson(data)
}

func (self *RESTFitting) delete(w rest.ResponseWriter, req *rest.Request) {
	metric := metric.MetricID(req.PathParam("metric"))

	self.writer.DeleteMetric(metric)

}
