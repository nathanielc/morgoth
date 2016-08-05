package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/influxdata/kapacitor/udf"
	"github.com/influxdata/kapacitor/udf/agent"
	"github.com/influxdata/wlog"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/fingerprinters/jsdiv"
	"github.com/nathanielc/morgoth/fingerprinters/kstest"
	"github.com/nathanielc/morgoth/fingerprinters/sigma"
)

const (
	defaultMinSupport     = 0.05
	defaultErrorTolerance = 0.01
	defaultConsensus      = 0.5
)

var socket = flag.String("socket", "", "Optional listen socket. If set then Morgoth will run in UDF socket mode, otherwise it will expect communication over STDIN/STDOUT.")
var logLevel = flag.String("log-level", "info", "Default log level, one of debug, info, warn or error.")

func main() {
	// Parse flags
	flag.Parse()

	// Setup logging
	log.SetOutput(wlog.NewWriter(os.Stderr))
	if err := wlog.SetLevelFromName(*logLevel); err != nil {
		log.Fatal("E! ", err)
	}

	if *socket == "" {
		a := agent.New(os.Stdin, os.Stdout)
		h := newHandler(a)
		a.Handler = h

		log.Println("I! Starting agent using STDIN/STDOUT")
		a.Start()
		err := a.Wait()
		if err != nil {
			log.Fatal("E! ", err)
		}
		log.Println("I! Agent finished")
	} else {
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

		// Setup signal handler to stop Server on various signals
		s.StopOnSignals(os.Interrupt, syscall.SIGTERM)

		log.Println("I! Socket server listening on", addr.String())
		err = s.Serve()
		if err != nil {
			log.Fatal("E! ", err)
		}
		log.Println("I! Socket server stopped")
	}
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
	}()
}

type fingerprinterInfo struct {
	init    initFingerprinterFunc
	options *udf.OptionInfo
}

// Function that creates a new instance of a fingerprinter
type createFingerprinterFunc func() morgoth.Fingerprinter

// Init createFingerprinterFunc from udf.OptionValues
type initFingerprinterFunc func(opts []*udf.OptionValue) (createFingerprinterFunc, error)

var fingerprinters = map[string]fingerprinterInfo{
	"sigma": {
		options: &udf.OptionInfo{ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
		init: func(args []*udf.OptionValue) (createFingerprinterFunc, error) {
			deviations := args[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
			if deviations <= 0 {
				return nil, fmt.Errorf("sigma: deviations must be > 0, got %f", deviations)
			}
			return func() morgoth.Fingerprinter {
				return sigma.New(deviations)
			}, nil
		},
	},
	"kstest": {
		options: &udf.OptionInfo{ValueTypes: []udf.ValueType{udf.ValueType_INT}},
		init: func(args []*udf.OptionValue) (createFingerprinterFunc, error) {
			confidence := args[0].Value.(*udf.OptionValue_IntValue).IntValue
			if confidence < 0 || confidence > 5 {
				return nil, fmt.Errorf("kstest: confidence must be in range [0,5], got %d", confidence)
			}
			return func() morgoth.Fingerprinter {
				return kstest.New(uint(confidence))
			}, nil
		},
	},
	"jsdiv": {
		options: &udf.OptionInfo{ValueTypes: []udf.ValueType{
			udf.ValueType_DOUBLE,
			udf.ValueType_DOUBLE,
			udf.ValueType_DOUBLE,
			udf.ValueType_DOUBLE,
		}},
		init: func(args []*udf.OptionValue) (createFingerprinterFunc, error) {
			min := args[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
			max := args[1].Value.(*udf.OptionValue_DoubleValue).DoubleValue
			binWidth := args[2].Value.(*udf.OptionValue_DoubleValue).DoubleValue
			pValue := args[3].Value.(*udf.OptionValue_DoubleValue).DoubleValue

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
	field          string
	scoreField     string
	minSupport     float64
	errorTolerance float64
	consensus      float64
	agent          *agent.Agent

	currentWindow *morgoth.Window
	beginBatch    *udf.BeginBatch
	batchPoints   []*udf.Point
	detectors     map[string]*morgoth.Detector

	fingerprinters []createFingerprinterFunc
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

// Return the InfoResponse. Describing the properties of this Handler
func (h *Handler) Info() (*udf.InfoResponse, error) {
	options := map[string]*udf.OptionInfo{
		"field":          {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
		"scoreField":     {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
		"minSupport":     {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
		"errorTolerance": {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
		"consensus":      {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
		"logLevel":       {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
	}
	// Add in options from fingerprinters
	for name, info := range fingerprinters {
		options[name] = info.options
	}
	info := &udf.InfoResponse{
		Wants:    udf.EdgeType_BATCH,
		Provides: udf.EdgeType_BATCH,
		Options:  options,
	}
	return info, nil

}

// Initialize the Handler with the provided options.
func (h *Handler) Init(r *udf.InitRequest) (*udf.InitResponse, error) {
	init := &udf.InitResponse{
		Success: true,
	}
	var errors []string
	for _, opt := range r.Options {
		switch opt.Name {
		case "field":
			h.field = opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
		case "scoreField":
			h.scoreField = opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
		case "minSupport":
			h.minSupport = opt.Values[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
		case "errorTolerance":
			h.errorTolerance = opt.Values[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
		case "consensus":
			h.consensus = opt.Values[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
		case "logLevel":
			level := opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
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
					h.fingerprinters = append(h.fingerprinters, createFn)
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
func (h *Handler) Snaphost() (*udf.SnapshotResponse, error) {
	return &udf.SnapshotResponse{}, nil
}

// Restore a previous snapshot.
func (h *Handler) Restore(*udf.RestoreRequest) (*udf.RestoreResponse, error) {
	return &udf.RestoreResponse{}, nil
}

// A batch has begun.
func (h *Handler) BeginBatch(b *udf.BeginBatch) error {
	h.currentWindow = &morgoth.Window{}
	h.beginBatch = b
	h.batchPoints = h.batchPoints[0:0]
	return nil
}

// A point has arrived.
func (h *Handler) Point(p *udf.Point) error {
	// Keep point around
	h.batchPoints = append(h.batchPoints, p)
	var value float64
	if f, ok := p.FieldsDouble[h.field]; ok {
		value = f
	} else {
		if i, ok := p.FieldsInt[h.field]; ok {
			value = float64(i)
		} else {
			return fmt.Errorf("no field %s is not a float or int", h.field)
		}
	}
	h.currentWindow.Data = append(h.currentWindow.Data, value)
	return nil
}

// The batch is complete.
func (h *Handler) EndBatch(b *udf.EndBatch) error {
	detector, ok := h.detectors[b.Group]
	if !ok {
		// We validated the args ourselves, ignore the error here
		detector, _ = morgoth.NewDetector(
			h.consensus,
			h.minSupport,
			h.errorTolerance,
			h.newFingerprinters(),
		)
		h.detectors[b.Group] = detector
	}
	if anomalous, avgSupport := detector.IsAnomalous(h.currentWindow); anomalous {
		// Send batch back to Kapacitor since it was anomalous
		h.agent.Responses <- &udf.Response{
			Message: &udf.Response_Begin{
				Begin: h.beginBatch,
			},
		}
		for _, p := range h.batchPoints {
			if h.scoreField != "" {
				p.FieldsDouble[h.scoreField] = 1 - avgSupport
			}
			h.agent.Responses <- &udf.Response{
				Message: &udf.Response_Point{
					Point: p,
				},
			}
		}
		h.agent.Responses <- &udf.Response{
			Message: &udf.Response_End{
				End: b,
			},
		}
	}
	return nil
}

// Gracefully stop the Handler.
// No other methods will be called.
func (h *Handler) Stop() {
	close(h.agent.Responses)
}

func (h *Handler) newFingerprinters() []morgoth.Fingerprinter {
	f := make([]morgoth.Fingerprinter, len(h.fingerprinters))
	for i, create := range h.fingerprinters {
		f[i] = create()
	}
	return f
}
