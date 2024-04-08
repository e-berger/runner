package results

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudwatchtype "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/prometheus/client_golang/prometheus"
)

type resultsPing struct {
	results
	statusCode string
}

type ResultsPingJSON struct {
	Id         string `json:"probe_id"`
	Time       string `json:"time"`
	Location   string `json:"location"`
	Latency    string `json:"latency"`
	ErrorProbe string `json:"error"`
	StatusCode string `json:"statusCode"`
}

func (r *resultsPing) String() string {
	return fmt.Sprintf("id: %s, location: %s, latency: %f, statusCode: %s",
		r.id, r.location, r.latency, r.statusCode)
}

func (r *resultsPing) GetId() string {
	return r.id
}

func (r *resultsPing) GetLatency() float64 {
	return r.latency
}

func (r *resultsPing) GetTime() time.Time {
	return r.time
}

func NewResultsPing(id string, location types.Location, started time.Time, latency int64, statusCode int) *resultsPing {
	return &resultsPing{
		results: results{
			id:       id,
			time:     started,
			location: location,
			latency:  float64(latency) / 1000.0,
		},
		statusCode: strconv.Itoa(statusCode),
	}
}

func NewResultsPingEmpty(id string, location types.Location) *resultsPing {
	return &resultsPing{
		results: results{
			id:       id,
			location: location,
			time:     time.Now(),
		},
	}
}

func (r *resultsPing) GetPrometheusMetrics() prometheus.Collector {
	completionTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sheepdog_runner",
		Name:      "request_duration_seconds",
		Help:      "Duration of the request.",
		Buckets:   []float64{0.1, 0.2, 0.3},
	}, []string{"method", "status", "location"})

	completionTime.With(prometheus.Labels{
		"status":   r.statusCode,
		"location": r.location.String(),
	}).Observe(r.latency)
	return completionTime
}

func (r *resultsPing) GetCloudWatchDimensions() []cloudwatchtype.Dimension {
	return []cloudwatchtype.Dimension{
		{
			Name:  aws.String("location"),
			Value: aws.String(r.location.String()),
		},
		{
			Name:  aws.String("status_code"),
			Value: aws.String(r.statusCode),
		},
	}
}

func (r *resultsPing) SetError(err error) {
	r.errorProbe = err
}

func (r *resultsPing) SetLatency(latency int64) {
	r.latency = float64(latency) / 1000.0
}

func (r *resultsPing) SetStatusCode(statusCode string) {
	r.statusCode = statusCode
}

func (r *resultsPing) SetTime(time time.Time) {
	r.time = time
}

func (r *resultsPing) GetErrorProbe() error {
	return r.errorProbe
}

func (r *resultsPing) MarshalJSON() ([]byte, error) {
	result := &ResultsHttpJSON{
		Id:         r.id,
		Time:       r.time.Format(time.RFC3339),
		Location:   r.location.String(),
		Latency:    fmt.Sprintf("%f", r.latency),
		ErrorProbe: fmt.Sprintf("%v", r.errorProbe),
		StatusCode: r.statusCode,
	}
	return json.Marshal(result)
}
