package collector

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cursorTimeoutMillis = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "parameters",
		Name:      "cursor_timeout_millis",
		Help:      "An integer that represents the cursorTimoutMillis option in mongod",
	})
	metric_mapping = map[string]prometheus.Gauge{
		"cursorTimeoutMillis": cursorTimeoutMillis,
	}
)

type ParameterMetrics struct {
}

func (p *ParameterMetrics) Export(ch chan<- prometheus.Metric) {
	cursorTimeoutMillis.Collect(ch)
}

func (p *ParameterMetrics) Describe(ch chan<- *prometheus.Desc) {
	cursorTimeoutMillis.Describe(ch)
}

func GetParameters(session *mgo.Session) *ParameterMetrics {
	for parameter, metric := range metric_mapping {
		result := make(map[string]interface{})
		err := session.DB("admin").Run(bson.D{{"getParameter", 1}, {parameter, 1}}, result)
		if err != nil {
			glog.Error("Failed to get parameter value for %v: %v", parameter, err)
			continue
		}
		if val, ok := result[parameter]; ok {
			metric.Set(float64(val.(int)))
		} else {
			glog.Error("Unexpected response from getParameter command: %v", result)
		}
	}
	return &ParameterMetrics{}
}
