package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudwatchtype "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/prometheus/client_golang/prometheus"
)

type metricsHttp struct {
	metrics
	httpMethod string
	statusCode string
}

func (m metricsHttp) String() string {
	return fmt.Sprintf("id: %s, location: %s, latency: %f, valid: %s, httpMethod: %s, statusCode: %s",
		m.id, m.location, m.latency, m.valid, m.httpMethod, m.statusCode)
}

func (m metricsHttp) GetId() string {
	return m.id
}

func (m metricsHttp) GetLatency() float64 {
	return m.latency
}

func (m metricsHttp) GetTime() time.Time {
	return m.time
}

func NewResultHttpDetails(id string, location types.Location, started time.Time, latency int64, valid int, httpMethod string, statusCode int) metricsHttp {
	return metricsHttp{
		metrics: metrics{
			id:       id,
			time:     started,
			location: location,
			latency:  float64(latency) / 1000.0,
			valid:    strconv.Itoa(valid),
		},
		httpMethod: httpMethod,
		statusCode: strconv.Itoa(statusCode),
	}
}

func (m metricsHttp) GetPrometheusMetrics() prometheus.Collector {
	completionTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sheepdog_runner",
		Name:      "request_duration_seconds",
		Help:      "Duration of the request.",
		Buckets:   []float64{0.1, 0.2, 0.3},
	}, []string{"method", "status", "location"})

	completionTime.With(prometheus.Labels{
		"method":   m.httpMethod,
		"status":   m.statusCode,
		"location": m.location.String(),
	}).Observe(m.latency)
	return completionTime
}

func (m metricsHttp) GetCloudWatchDimensions() []cloudwatchtype.Dimension {
	return []cloudwatchtype.Dimension{
		{
			Name:  aws.String("location"),
			Value: aws.String(m.location.String()),
		},
		{
			Name:  aws.String("status_code"),
			Value: aws.String(m.statusCode),
		},
	}
}
