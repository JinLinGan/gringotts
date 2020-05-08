package checks

import (
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/global"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
	"github.com/jinlingan/gringotts/gringotts-agent/telegraf"
	_ "github.com/jinlingan/gringotts/gringotts-agent/telegraf/inputs/all"
)

type Initializer interface {
	// Init performs one time setup of the plugin and returns an error if the
	// configuration is invalid.
	Init() error
}

type TelegrafCheck struct {
	JobID string

	TelegrafInput telegraf.Input

	NameOverride      string
	MeasurementPrefix string
	MeasurementSuffix string
	Tags              map[string]string
	Filter            telegraf.Filter
}

func (t *TelegrafCheck) LogName() string {
	return t.JobID
}

func (t *TelegrafCheck) Run() error {
	//panic("implement me")

	sender := global.GlobalSenderPool.GetSender(t.JobID)

	acc := telegraf.NewAccumulator(t, sender)
	acc.SetPrecision(time.Nanosecond)

	return t.TelegrafInput.Gather(acc)
}

func (t *TelegrafCheck) Stop() {
	//panic("implement me")
}

func (t *TelegrafCheck) Init() error {
	if p, ok := t.TelegrafInput.(Initializer); ok {
		err := p.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TelegrafCheck) metricFiltered(metric model.Metric) {
	metric.Drop()
}
func (t *TelegrafCheck) MakeMetric(metric model.Metric) model.Metric {
	if ok := t.Filter.Select(metric); !ok {
		t.metricFiltered(metric)
		return nil
	}

	m := makemetric(
		metric,
		t.NameOverride,
		t.MeasurementPrefix,
		t.MeasurementSuffix,
		t.Tags)

	t.Filter.Modify(metric)
	if len(metric.FieldList()) == 0 {
		t.metricFiltered(metric)
		return nil
	}

	return m
}

func makemetric(
	metric model.Metric,
	nameOverride string,
	namePrefix string,
	nameSuffix string,
	tags map[string]string,
) model.Metric {
	if len(nameOverride) != 0 {
		metric.SetName(nameOverride)
	}

	if len(namePrefix) != 0 {
		metric.AddPrefix(namePrefix)
	}
	if len(nameSuffix) != 0 {
		metric.AddSuffix(nameSuffix)
	}

	// Apply plugin-wide tags
	for k, v := range tags {
		if _, ok := metric.GetTag(k); !ok {
			metric.AddTag(k, v)
		}
	}

	return metric
}
