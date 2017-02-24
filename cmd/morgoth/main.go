package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/influxdata/kapacitor/udf/agent"
	"github.com/influxdata/wlog"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/counter"
	"github.com/nathanielc/morgoth/fingerprinters/jsdiv"
	"github.com/nathanielc/morgoth/fingerprinters/kstest"
	"github.com/nathanielc/morgoth/fingerprinters/sigma"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultMinSupport      = 0.05
	defaultErrorTolerance  = 0.01
	defaultConsensus       = 0.5
	defaultMetricsBindAddr = ":6767"
)

var socket = flag.String("socket", "", "Optional listen socket. If set then Morgoth will run in UDF socket mode, otherwise it will expect communication over STDIN/STDOUT.")
var logLevel = flag.String("log-level", "info", "Default log level, one of debug, info, warn or error.")
var metricsBind = flag.String("metrics-bind", defaultMetricsBindAddr, "Bind address of the metrics HTTP server. The metrics server will only start if also using the socket mode of operation.")

var detectorGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "morgoth_detectors",
	Help: "Current number of active detectors.",
})

func init() {
	prometheus.MustRegister(detectorGauge)
}

var detectorDims = []string{
	"task",
	"node",
	"group",
}

var fingerprinterDims = append(detectorDims, "fingerprinter")

var windowCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "morgoth_windows_total",
		Help: "Number of windows processed.",
	},
	detectorDims,
)
var pointCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "morgoth_points_total",
		Help: "Number of points processed.",
	},
	detectorDims,
)
var anomalousCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "morgoth_anomalies_total",
		Help: "Number of anomalies detected.",
	},
	detectorDims,
)
var uniqueFingerpintsGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "morgoth_unique_fingerprints",
		Help: "Current number of unique fingerprints.",
	},
	fingerprinterDims,
)

func main() {
	// Parse flags
	flag.Parse()

	// Setup logging
	log.SetOutput(wlog.NewWriter(os.Stderr))
	if err := wlog.SetLevelFromName(*logLevel); err != nil {
		log.Fatal("E! ", err)
	}

	// Create error channels
	metricsErr := make(chan error, 1)
	errC := make(chan error, 1)
	if *socket == "" {
		a := agent.New(os.Stdin, os.Stdout)
		h := newHandler(a)
		a.Handler = h
		defer h.Stop()

		log.Println("I! Starting agent using STDIN/STDOUT")
		a.Start()

		go func() {
			errC <- a.Wait()
		}()
	} else {
		// Start the metrics server.
		// Only start the metrics server in socket mode or the bind address would conflict with each new process.
		go func() {
			log.Println("I! Starting metrics HTTP server on", *metricsBind)
			http.Handle("/metrics", promhttp.Handler())
			metricsErr <- http.ListenAndServe(*metricsBind, nil)
		}()

		// Create unix socket
		addr, err := net.ResolveUnixAddr("unix", *socket)
		if err != nil {
			log.Fatal("E! ", err)
		}
		l, err := net.ListenUnix("unix", addr)
		if err != nil {
			log.Fatal("E! ", err)
		}

		// Create server that listens on the socket
		s := agent.NewServer(l, &accepter{})
		defer s.Stop() // this closes the listener

		// Setup signal handler to stop Server on various signals
		s.StopOnSignals(os.Interrupt, os.Kill)

		go func() {
			log.Println("I! Starting socket server on", addr.String())
			errC <- s.Serve()
		}()
	}

	select {
	case err := <-metricsErr:
		if err != nil {
			log.Println("E!", err)
		}
	case err := <-errC:
		if err != nil {
			log.Println("E!", err)
		}
	}
	log.Println("I! Stopping")
}

// Simple connection accepter
type accepter struct {
	count int64
}

// Create a new agent/handler for each new connection.
// Count and log each new connection and termination.
func (acc *accepter) Accept(conn net.Conn) {
	count := acc.count
	acc.count++
	a := agent.New(conn, conn)
	h := newHandler(a)
	a.Handler = h

	log.Println("I! Starting agent for connection", count)
	a.Start()
	go func() {
		err := a.Wait()
		if err != nil {
			log.Printf("E! Agent for connection %d terminated with error: %s", count, err)
		} else {
			log.Printf("I! Agent for connection %d finished", count)
		}
		h.Close()
	}()
}

type fingerprinterInfo struct {
	init    initFingerprinterFunc
	options *agent.OptionInfo
}

// Function that creates a new instance of a fingerprinter
type createFingerprinterFunc func() morgoth.Fingerprinter

// Init createFingerprinterFunc from agent.OptionValues
type initFingerprinterFunc func(opts []*agent.OptionValue) (createFingerprinterFunc, error)

var fingerprinters = map[string]fingerprinterInfo{
	"sigma": {
		options: &agent.OptionInfo{ValueTypes: []agent.ValueType{agent.ValueType_DOUBLE}},
		init: func(args []*agent.OptionValue) (createFingerprinterFunc, error) {
			deviations := args[0].Value.(*agent.OptionValue_DoubleValue).DoubleValue
			if deviations <= 0 {
				return nil, fmt.Errorf("sigma: deviations must be > 0, got %f", deviations)
			}
			return func() morgoth.Fingerprinter {
				return sigma.New(deviations)
			}, nil
		},
	},
	"kstest": {
		options: &agent.OptionInfo{ValueTypes: []agent.ValueType{agent.ValueType_INT}},
		init: func(args []*agent.OptionValue) (createFingerprinterFunc, error) {
			confidence := args[0].Value.(*agent.OptionValue_IntValue).IntValue
			if confidence < 0 || confidence > 5 {
				return nil, fmt.Errorf("kstest: confidence must be in range [0,5], got %d", confidence)
			}
			return func() morgoth.Fingerprinter {
				return kstest.New(uint(confidence))
			}, nil
		},
	},
	"jsdiv": {
		options: &agent.OptionInfo{ValueTypes: []agent.ValueType{
			agent.ValueType_DOUBLE,
			agent.ValueType_DOUBLE,
			agent.ValueType_DOUBLE,
			agent.ValueType_DOUBLE,
		}},
		init: func(args []*agent.OptionValue) (createFingerprinterFunc, error) {
			min := args[0].Value.(*agent.OptionValue_DoubleValue).DoubleValue
			max := args[1].Value.(*agent.OptionValue_DoubleValue).DoubleValue
			binWidth := args[2].Value.(*agent.OptionValue_DoubleValue).DoubleValue
			pValue := args[3].Value.(*agent.OptionValue_DoubleValue).DoubleValue

			if binWidth <= 0 {
				return nil, fmt.Errorf("jsdiv: binWidth, arg 3, must be > 0, got %f", binWidth)
			}
			if pValue <= 0 || pValue > 1 {
				return nil, fmt.Errorf("jsdiv: pValue, arg 4, must be in range (0,1], got %f", pValue)
			}
			if (max-min)/binWidth < 3 {
				return nil, fmt.Errorf("jsdiv: more than 3 bins should fit in the range [min,max]")
			}

			return func() morgoth.Fingerprinter {
				return jsdiv.New(min, max, binWidth, pValue)
			}, nil
		},
	},
}

// A Kapacitor UDF Handler
type Handler struct {
	taskID string
	nodeID string

	field          string
	scoreField     string
	minSupport     float64
	errorTolerance float64
	consensus      float64
	agent          *agent.Agent

	currentWindow *morgoth.Window
	beginBatch    *agent.BeginBatch
	batchPoints   []*agent.Point
	detectors     map[string]*morgoth.Detector

	fingerprinters []fingerprinterCreator
}

type fingerprinterCreator struct {
	Kind   string
	Create createFingerprinterFunc
}

func newHandler(a *agent.Agent) *Handler {
	return &Handler{
		agent:          a,
		minSupport:     defaultMinSupport,
		errorTolerance: defaultErrorTolerance,
		consensus:      defaultConsensus,
		detectors:      make(map[string]*morgoth.Detector),
	}
}

func (h *Handler) Close() {
	for _, d := range h.detectors {
		d.Close()
	}
}

func (h *Handler) detectorName(group string) string {
	return fmt.Sprintf("%s:%s,group=%s", h.taskID, h.nodeID, group)
}

// Return the InfoResponse. Describing the properties of this Handler
func (h *Handler) Info() (*agent.InfoResponse, error) {
	options := map[string]*agent.OptionInfo{
		"field":          {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
		"scoreField":     {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
		"minSupport":     {ValueTypes: []agent.ValueType{agent.ValueType_DOUBLE}},
		"errorTolerance": {ValueTypes: []agent.ValueType{agent.ValueType_DOUBLE}},
		"consensus":      {ValueTypes: []agent.ValueType{agent.ValueType_DOUBLE}},
		"logLevel":       {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
	}
	// Add in options from fingerprinters
	for name, info := range fingerprinters {
		options[name] = info.options
	}
	info := &agent.InfoResponse{
		Wants:    agent.EdgeType_BATCH,
		Provides: agent.EdgeType_BATCH,
		Options:  options,
	}
	return info, nil

}

// Initialize the Handler with the provided options.
func (h *Handler) Init(r *agent.InitRequest) (*agent.InitResponse, error) {
	h.taskID = r.TaskID
	h.nodeID = r.NodeID

	init := &agent.InitResponse{
		Success: true,
	}
	var errors []string
	for _, opt := range r.Options {
		switch opt.Name {
		case "field":
			h.field = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
		case "scoreField":
			h.scoreField = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
		case "minSupport":
			h.minSupport = opt.Values[0].Value.(*agent.OptionValue_DoubleValue).DoubleValue
		case "errorTolerance":
			h.errorTolerance = opt.Values[0].Value.(*agent.OptionValue_DoubleValue).DoubleValue
		case "consensus":
			h.consensus = opt.Values[0].Value.(*agent.OptionValue_DoubleValue).DoubleValue
		case "logLevel":
			level := opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
			err := wlog.SetLevelFromName(level)
			if err != nil {
				init.Success = false
				errors = append(errors, err.Error())
			}
		default:
			if info, ok := fingerprinters[opt.Name]; ok {
				createFn, err := info.init(opt.Values)
				if err != nil {
					init.Success = false
					errors = append(errors, err.Error())
				} else {
					h.fingerprinters = append(h.fingerprinters, fingerprinterCreator{
						Kind:   opt.Name,
						Create: createFn,
					})
				}
			} else {
				return nil, fmt.Errorf("received unknown init option %q", opt.Name)
			}
		}
	}

	if h.field == "" {
		errors = append(errors, "field must not be empty")
	}
	if h.minSupport < 0 || h.minSupport > 1 {
		errors = append(errors, "minSupport must be in the range [0,1)")
	}
	if h.errorTolerance < 0 || h.errorTolerance > 1 {
		errors = append(errors, "errorTolerance must be in the range [0,1)")
	}
	if (h.consensus != -1 && h.consensus < 0) || h.consensus > 1 {
		errors = append(errors, "consensus must be in the range [0,1) or equal to -1")
	}
	if h.minSupport <= h.errorTolerance {
		errors = append(errors, "invalid minSupport or errorTolerance: minSupport must be greater than errorTolerance")
	}
	init.Success = len(errors) == 0
	init.Error = strings.Join(errors, "\n")

	log.Printf("D! %#v", h)

	return init, nil
}

// Create a snapshot of the running state of the handler.
func (h *Handler) Snapshot() (*agent.SnapshotResponse, error) {
	return &agent.SnapshotResponse{}, nil
}

// Restore a previous snapshot.
func (h *Handler) Restore(*agent.RestoreRequest) (*agent.RestoreResponse, error) {
	return &agent.RestoreResponse{}, nil
}

// A batch has begun.
func (h *Handler) BeginBatch(b *agent.BeginBatch) error {
	h.currentWindow = &morgoth.Window{}
	h.beginBatch = b
	h.batchPoints = h.batchPoints[0:0]
	return nil
}

// A point has arrived.
func (h *Handler) Point(p *agent.Point) error {
	// Keep point around
	h.batchPoints = append(h.batchPoints, p)
	var value float64
	if f, ok := p.FieldsDouble[h.field]; ok {
		value = f
	} else {
		if i, ok := p.FieldsInt[h.field]; ok {
			value = float64(i)
		} else {
			return fmt.Errorf("field %q is not a float or int", h.field)
		}
	}
	h.currentWindow.Data = append(h.currentWindow.Data, value)
	return nil
}

// The batch is complete.
func (h *Handler) EndBatch(b *agent.EndBatch) error {
	detector, ok := h.detectors[b.Group]
	if !ok {
		metrics := h.createDetectorMetrics(b.Group)
		if err := metrics.Register(); err != nil {
			return errors.Wrapf(err, "failed to register metrics for group: %q", b.Group)
		}
		// We validated the args ourselves, ignore the error here
		detector, _ = morgoth.NewDetector(
			metrics,
			h.consensus,
			h.minSupport,
			h.errorTolerance,
			h.newFingerprinters(),
		)
		h.detectors[b.Group] = detector
		detectorGauge.Inc()
	}
	if anomalous, avgSupport := detector.IsAnomalous(h.currentWindow); anomalous {
		// Send batch back to Kapacitor since it was anomalous
		h.agent.Responses <- &agent.Response{
			Message: &agent.Response_Begin{
				Begin: h.beginBatch,
			},
		}
		for _, p := range h.batchPoints {
			if h.scoreField != "" {
				if p.FieldsDouble == nil {
					p.FieldsDouble = make(map[string]float64, 1)
				}
				p.FieldsDouble[h.scoreField] = 1 - avgSupport
			}
			h.agent.Responses <- &agent.Response{
				Message: &agent.Response_Point{
					Point: p,
				},
			}
		}
		h.agent.Responses <- &agent.Response{
			Message: &agent.Response_End{
				End: b,
			},
		}
	}
	return nil
}

func (h *Handler) createDetectorMetrics(group string) *morgoth.DetectorMetrics {
	labels := prometheus.Labels{
		"task":  h.taskID,
		"node":  h.nodeID,
		"group": group,
	}
	metrics := &morgoth.DetectorMetrics{
		WindowCount:          windowCount.With(labels),
		PointCount:           pointCount.With(labels),
		AnomalousCount:       anomalousCount.With(labels),
		FingerprinterMetrics: make([]*counter.Metrics, len(h.fingerprinters)),
	}
	for i, creator := range h.fingerprinters {
		labels = prometheus.Labels{
			"task":          h.taskID,
			"node":          h.nodeID,
			"group":         group,
			"fingerprinter": fmt.Sprintf("%s-%d", creator.Kind, i),
		}
		metrics.FingerprinterMetrics[i] = &counter.Metrics{
			UniqueFingerprints: uniqueFingerpintsGauge.With(labels),
		}
	}
	return metrics
}

// Gracefully stop the Handler.
// No other methods will be called.
func (h *Handler) Stop() {
	close(h.agent.Responses)
}

func (h *Handler) newFingerprinters() []morgoth.Fingerprinter {
	f := make([]morgoth.Fingerprinter, len(h.fingerprinters))
	for i, creator := range h.fingerprinters {
		f[i] = creator.Create()
	}
	return f
}
