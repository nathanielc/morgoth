package main

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/kapacitor/udf/agent"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/counter"
)

func main() {
	a := agent.New()
	h := newHandler()
	a.Handler = h

	log.Println("Starting agent")
	a.Start()
	err := a.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

// A Kapacitor UDF Handler
type Handler struct {
	field          string
	minSupport     float64
	errorTolerance float64
	normalCount    int
	consensus      float64

	currentBatch *morgoth.Window
	detectors    map[string]*morgoth.Detector
}

// Pair of fingerprinter and counter
type fingerprinterCounter struct {
	fingerprinter Fingerprinter
	counter       counter.Counter
}

func newHandler() *Handler {
	return &Handler{}
}

// Return the InfoResponse. Describing the properties of this Handler
func (h *Handler) Info() (*udf.InfoResponse, error) {
	info := &udf.InfoResponse{
		Wants:    udf.EdgeType_BATCH,
		Provides: udf.EdgeType_STREAM,
		Options: map[string]*udf.OptionInfo{
			"field":          {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
			"minSupport":     {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
			"errorTolerance": {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
			"normalCount":    {ValueTypes: []udf.ValueType{udf.ValueType_INT}},
			"consensus":      {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
			// Fingerprinters
			"sigma":  {ValueTypes: []udf.ValueType{udf.ValueType_DOUBLE}},
			"kstest": {ValueTypes: []udf.ValueType{udf.ValueType_INT}},
			"jsdiv": {ValueTypes: []udf.ValueType{
				udf.ValueType_DOUBLE,
				udf.ValueType_DOUBLE,
				udf.ValueType_INT,
				udf.ValueType_DOUBLE,
			}},
		},
	}
	return info, nil

}

// Initialize the Handler with the provided options.
func (h *Handler) Init(*udf.InitRequest) (*udf.InitResponse, error) {
	init := &udf.InitResponse{
		Success: true,
		Error:   "",
	}
	for _, opt := range r.Options {
		switch opt.Name {
		case "field":
			h.field = opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
		case "minSupport":
			h.minSupport = opt.Values[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
		case "errorTolerance":
			h.errorTolerance = opt.Values[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
		case "normalCount":
			h.normalCount = int(opt.Values[0].Value.(*udf.OptionValue_IntValue).IntValue)
		case "consensus":
			h.consensus = opt.Values[0].Value.(*udf.OptionValue_DoubleValue).DoubleValue
		}
	}

	if h.field == "" {
		init.Success = false
		init.Error += "field must not be empty"
	}
	if h.minSupport < 0 || h.minSupport > 1 {
		init.Success = false
		init.Error += "minSupport must be in the range [0,1) "
	}
	if h.errorTolerance < 0 || h.errorTolerance > 1 {
		init.Success = false
		init.Error += "errorTolerance must be in the range [0,1) "
	}
	if h.normalCount < 0 {
		init.Success = false
		init.Error += "normalCount must be greater than 0 "
	}
	if h.consensus < 0 || h.consensus > 1 {
		init.Success = false
		init.Error += "consensus must be in the range [0,1) "
	}

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
	h.currentBatch = &morgoth.Window{
		Name: b.Name,
		Tags: b.Tags,
	}
	return nil
}

// A point has arrived.
func (h *Handler) Point(p *udf.Point) error {
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
	if h.currentBatch.Start.IsZero() {
		h.currentBatch.Start = time.Unix(0, p.Time).UTC()
	}
	h.currentBatch.Data = append(h.currentBatch, value)
	return nil
}

// The batch is complete.
func (h *Handler) EndBatch(b *udf.EndBatch) error {
	h.currentBatch.Stop = time.Unix(0, b.TMax).UTC()
	detector, ok := h.detectors[b.Group]
	if !ok {
		detector = morgoth.NewDetector(
			h.normalCount,
			h.consensus,
			h.minSupport,
			h.errorTolerance,
			h.newFingerprinters(),
		)
		h.detectors[b.Group] = detector
	}
	if detector.IsAnomalous(h.currentBatch) {
		// Send point back to Kapacitor
	}

	return nil
}

// Gracefully stop the Handler.
// No other methods will be called.
func (h *Handler) Stop() {}

func (h *Handler) newFingerprinters() []morgoth.Fingerprinter {
	return nil
}
