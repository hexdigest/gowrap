package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/prometheus template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i AnotherTestInterface -t ../templates/prometheus -o interface_with_prometheus_metric_prefix.go -v MetricPrefix=some_prefix

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// AnotherTestInterfaceWithPrometheus implements AnotherTestInterface interface with all methods wrapped
// with Prometheus metrics
type AnotherTestInterfaceWithPrometheus struct {
	base         AnotherTestInterface
	instanceName string
}

var anothertestinterfaceDurationSummaryVec = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "some_prefix_duration_seconds",
		Help:       "some_prefix runtime duration and result",
		MaxAge:     time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"instance_name", "method", "result"})

// NewAnotherTestInterfaceWithPrometheus returns an instance of the AnotherTestInterface decorated with prometheus summary metric
func NewAnotherTestInterfaceWithPrometheus(base AnotherTestInterface, instanceName string) AnotherTestInterfaceWithPrometheus {
	return AnotherTestInterfaceWithPrometheus{
		base:         base,
		instanceName: instanceName,
	}
}

// Channels implements AnotherTestInterface
func (_d AnotherTestInterfaceWithPrometheus) Channels(chA chan bool, chB chan<- bool, chanC <-chan bool) {
	_since := time.Now()
	defer func() {
		result := "ok"
		anothertestinterfaceDurationSummaryVec.WithLabelValues(_d.instanceName, "Channels", result).Observe(time.Since(_since).Seconds())
	}()
	_d.base.Channels(chA, chB, chanC)
	return
}

// F implements AnotherTestInterface
func (_d AnotherTestInterfaceWithPrometheus) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		anothertestinterfaceDurationSummaryVec.WithLabelValues(_d.instanceName, "F", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.F(ctx, a1, a2...)
}

// NoError implements AnotherTestInterface
func (_d AnotherTestInterfaceWithPrometheus) NoError(s1 string) (s2 string) {
	_since := time.Now()
	defer func() {
		result := "ok"
		anothertestinterfaceDurationSummaryVec.WithLabelValues(_d.instanceName, "NoError", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.NoError(s1)
}

// NoParamsOrResults implements AnotherTestInterface
func (_d AnotherTestInterfaceWithPrometheus) NoParamsOrResults() {
	_since := time.Now()
	defer func() {
		result := "ok"
		anothertestinterfaceDurationSummaryVec.WithLabelValues(_d.instanceName, "NoParamsOrResults", result).Observe(time.Since(_since).Seconds())
	}()
	_d.base.NoParamsOrResults()
	return
}
