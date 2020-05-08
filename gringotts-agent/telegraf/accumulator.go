package telegraf

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/global"

	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

// Accumulator allows adding metrics to the processing flow.
type Accumulator interface {
	// AddFields adds a metric to the accumulator with the given measurement
	// name, fields, and tags (and timestamp). If a timestamp is not provided,
	// then the accumulator sets it to "now".
	AddFields(measurement string,
		fields map[string]interface{},
		tags map[string]string,
		t ...time.Time)

	// AddGauge is the same as AddFields, but will add the metric as a "Gauge" type
	AddGauge(measurement string,
		fields map[string]interface{},
		tags map[string]string,
		t ...time.Time)

	// AddCounter is the same as AddFields, but will add the metric as a "Counter" type
	AddCounter(measurement string,
		fields map[string]interface{},
		tags map[string]string,
		t ...time.Time)

	// AddSummary is the same as AddFields, but will add the metric as a "Summary" type
	AddSummary(measurement string,
		fields map[string]interface{},
		tags map[string]string,
		t ...time.Time)

	// AddHistogram is the same as AddFields, but will add the metric as a "Histogram" type
	AddHistogram(measurement string,
		fields map[string]interface{},
		tags map[string]string,
		t ...time.Time)

	// AddMetric adds an metric to the accumulator.
	AddMetric(model.Metric)

	// SetPrecision sets the timestamp rounding precision.  All metrics addeds
	// added to the accumulator will have their timestamp rounded to the
	// nearest multiple of precision.
	SetPrecision(precision time.Duration)

	// Report an error.
	AddError(err error)

	// Upgrade to a TrackingAccumulator with space for maxTracked
	// metrics/batches.

	//TODO：暂时不实现这个
	//WithTracking(maxTracked int) TrackingAccumulator

	NFields() uint64
	NErrors() uint64
}

// TrackingID uniquely identifies a tracked metric group
type TrackingID uint64

// DeliveryInfo provides the results of a delivered metric group.
type DeliveryInfo interface {
	// ID is the TrackingID
	ID() TrackingID

	// Delivered returns true if the metric was processed successfully.
	Delivered() bool
}

// TrackingAccumulator is an Accumulator that provides a signal when the
// metric has been fully processed.  Sending more metrics than the accumulator
// has been allocated for without reading status from the Accepted or Rejected
// channels is an error.
type TrackingAccumulator interface {
	Accumulator

	// Add the Metric and arrange for tracking feedback after processing..
	AddTrackingMetric(m model.Metric) TrackingID

	// Add a group of Metrics and arrange for a signal when the group has been
	// processed.
	AddTrackingMetricGroup(group []model.Metric) TrackingID

	// Delivered returns a channel that will contain the tracking results.
	Delivered() <-chan DeliveryInfo
}

type MetricMaker interface {
	LogName() string
	MakeMetric(metric model.Metric) model.Metric
}

type accumulator struct {
	//sync.Mutex
	nFields   uint64
	nErrors   uint64
	maker     MetricMaker
	sender    global.Sender
	precision time.Duration
}

func NewAccumulator(
	maker MetricMaker,
	sender global.Sender,
) Accumulator {
	acc := accumulator{
		maker:     maker,
		sender:    sender,
		precision: time.Nanosecond,
	}
	return &acc
}

func (ac *accumulator) AddFields(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, model.Untyped, t...)
}

func (ac *accumulator) AddGauge(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, model.Gauge, t...)
}

func (ac *accumulator) AddCounter(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, model.Counter, t...)
}

func (ac *accumulator) AddSummary(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, model.Summary, t...)
}

func (ac *accumulator) AddHistogram(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, model.Histogram, t...)
}

func (ac *accumulator) AddMetric(m model.Metric) {
	m.SetTime(m.Time().Round(ac.precision))
	if m := ac.maker.MakeMetric(m); m != nil {
		//ac.metrics <- m
		ac.sender.SendTelegraf(m)
	}
}

func (ac *accumulator) addFields(
	measurement string,
	tags map[string]string,
	fields map[string]interface{},
	tp model.ValueType,
	t ...time.Time,
) {
	m, err := model.NewMetrics(measurement, tags, fields, ac.getTime(t), tp)
	if err != nil {
		return
	}
	if m := ac.maker.MakeMetric(m); m != nil {
		//ac.metrics <- m
		ac.sender.SendTelegraf(m)
		atomic.AddUint64(&ac.nFields, uint64(len(fields)))
	}
}

func (ac *accumulator) NFields() uint64 {
	return ac.nFields
}

func (ac *accumulator) NErrors() uint64 {
	return ac.nErrors
}

// AddError passes a runtime error to the accumulator.
// The error will be tagged with the plugin name and written to the log.
func (ac *accumulator) AddError(err error) {
	if err == nil {
		return
	}
	//NErrors.Incr(1)
	atomic.AddUint64(&ac.nErrors, 1)
	log.Printf("E! [%s] Error in plugin: %v", ac.maker.LogName(), err)
}

func (ac *accumulator) SetPrecision(precision time.Duration) {
	ac.precision = precision
}

func (ac *accumulator) getTime(t []time.Time) time.Time {
	var timestamp time.Time
	if len(t) > 0 {
		timestamp = t[0]
	} else {
		timestamp = time.Now()
	}
	return timestamp.Round(ac.precision)
}

//func (ac *accumulator) WithTracking(maxTracked int) telegraf.TrackingAccumulator {
//	return &trackingAccumulator{
//		Accumulator: ac,
//		delivered:   make(chan telegraf.DeliveryInfo, maxTracked),
//	}
//}
