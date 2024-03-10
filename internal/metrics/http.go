package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsHttp struct {
	Metrics
	HttpMethod string `json:"http_method"`
	StatusCode string `json:"status_code"`
}

func (r *MetricsHttp) GetId() string {
	return r.Id
}

func (r *MetricsHttp) GetLocation() string {
	return r.Location
}

func (r *MetricsHttp) GetLatency() float64 {
	return r.Latency
}

func (r *MetricsHttp) GeStatusCode() string {
	return r.StatusCode
}

func (r *MetricsHttp) String() string {
	return fmt.Sprintf("Id: %s, Location: %s, Latency: %f, Valid: %s, HttpMethod: %s, StatusCode: %s", r.Id, r.Location, r.Latency, r.Valid, r.HttpMethod, r.StatusCode)
}

func NewResultHttpDetails(id string, location int, started time.Time, latency int64, valid int, httpMethod string, statusCode int) *MetricsHttp {
	return &MetricsHttp{
		Metrics: Metrics{
			Id:       id,
			Time:     started,
			Location: strconv.Itoa(location),
			Latency:  float64(latency) / 1000.0,
			Valid:    strconv.Itoa(valid),
		},
		HttpMethod: httpMethod,
		StatusCode: strconv.Itoa(statusCode),
	}
}

func (m *MetricsHttp) GetMetrics() prometheus.Collector {
	completionTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sheepdog_runner",
		Name:      "request_duration_seconds",
		Help:      "Duration of the request.",
		Buckets:   []float64{0.1, 0.2, 0.3},
	}, []string{"method", "status", "location"})

	completionTime.With(prometheus.Labels{
		"method":   m.HttpMethod,
		"status":   m.StatusCode,
		"location": m.Location,
	}).Observe(m.Latency)
	return completionTime
}

func (m *MetricsHttp) GetTime() time.Time {
	return m.Time
}
