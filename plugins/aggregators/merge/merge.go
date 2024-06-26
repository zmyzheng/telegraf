//go:generate ../../../tools/readme_config_includer/generator
package merge

import (
	_ "embed"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/aggregators"
)

//go:embed sample.conf
var sampleConfig string

type Merge struct {
	RoundTimestamp config.Duration `toml:"round_timestamp_to"`
	grouper        *metric.SeriesGrouper
}

func (*Merge) SampleConfig() string {
	return sampleConfig
}

func (a *Merge) Init() error {
	a.grouper = metric.NewSeriesGrouper()
	return nil
}

func (a *Merge) Add(m telegraf.Metric) {
	gm := m
	if a.RoundTimestamp > 0 {
		if unwrapped, ok := m.(telegraf.UnwrappableMetric); ok {
			gm = unwrapped.Unwrap().Copy()
		} else {
			gm = m.Copy()
		}
		ts := gm.Time()
		gm.SetTime(ts.Round(time.Duration(a.RoundTimestamp)))
	}
	a.grouper.AddMetric(gm)
}

func (a *Merge) Push(acc telegraf.Accumulator) {
	// Always use nanosecond precision to avoid rounding metrics that were
	// produced at a precision higher than the agent default.
	acc.SetPrecision(time.Nanosecond)

	for _, m := range a.grouper.Metrics() {
		acc.AddMetric(m)
	}
}

func (a *Merge) Reset() {
	a.grouper = metric.NewSeriesGrouper()
}

func init() {
	aggregators.Add("merge", func() telegraf.Aggregator {
		return &Merge{}
	})
}
