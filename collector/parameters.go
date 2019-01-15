package collector

import (
	"reflect"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricMapping = map[string]prometheus.Gauge{}
)

type ParameterMetrics struct {
}

func (p *ParameterMetrics) Export(ch chan<- prometheus.Metric) {
	for _, metric := range metricMapping {
		metric.Collect(ch)
	}
}

func (p *ParameterMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range metricMapping {
		metric.Describe(ch)
	}
}

func GetParameters(session *mgo.Session, parameters string) *ParameterMetrics {
	splitParameters := strings.Split(parameters, ",")
	for _, parameter := range splitParameters {
		if _, ok := metricMapping[parameter]; !ok {
			metricMapping[parameter] = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: "parameters",
				Name:      strings.ToLower(parameter),
				Help:      "A setParameter option in mongod",
			})
		}
		metric := metricMapping[parameter]
		result := make(map[string]interface{})
		err := session.DB("admin").Run(bson.D{{"getParameter", 1}, {parameter, 1}}, result)
		if err != nil {
			glog.Errorf("Failed to get parameter value for %v: %v", parameter, err)
			continue
		}
		if val, ok := result[parameter]; ok {
			switch valTyped := val.(type) {
			case int:
				metric.Set(float64(valTyped))
			case int32:
				metric.Set(float64(valTyped))
			case int64:
				metric.Set(float64(valTyped))
			case float64:
				metric.Set(valTyped)
			case bool:
				var bit int8
				if valTyped {
					bit = 1
				}
				metric.Set(float64(bit))
			default:
				glog.Errorf("Unknown parameter value for %v: %v of type %v", parameter, valTyped, reflect.TypeOf(val))
			}
		} else {
			glog.Errorf("Unexpected response from getParameter command: %v", result)
		}
	}
	return &ParameterMetrics{}
}
