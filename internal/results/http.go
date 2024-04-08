package results

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudwatchtype "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/prometheus/client_golang/prometheus"
)

type resultsHttp struct {
	results
	httpMethod string
	statusCode string
}

type ResultsHttpJSON struct {
	Id         string `json:"probe_id"`
	Time       string `json:"time"`
	Location   string `json:"location"`
	Latency    string `json:"latency"`
	ErrorProbe string `json:"error"`
	HttpMethod string `json:"httpMethod"`
	StatusCode string `json:"statusCode"`
}

func (r *resultsHttp) String() string {
	return fmt.Sprintf("id: %s, location: %s, latency: %f, httpMethod: %s, statusCode: %s",
		r.id, r.location, r.latency, r.httpMethod, r.statusCode)
}

func (r *resultsHttp) GetId() string {
	return r.id
}

func (r *resultsHttp) GetLatency() float64 {
	return r.latency
}

func (r *resultsHttp) GetTime() time.Time {
	return r.time
}

func NewResultsHttp(id string, location types.Location, started time.Time, latency int64, httpMethod string, statusCode int) *resultsHttp {
	return &resultsHttp{
		results: results{
			id:       id,
			time:     started,
			location: location,
			latency:  float64(latency) / 1000.0,
		},
		httpMethod: httpMethod,
		statusCode: strconv.Itoa(statusCode),
	}
}

func NewResultsHttpEmpty(id string, location types.Location, httpMethod string) *resultsHttp {
	return &resultsHttp{
		results: results{
			id:       id,
			location: location,
			time:     time.Now(),
		},
		httpMethod: httpMethod,
	}
}

func (r *resultsHttp) GetPrometheusMetrics() prometheus.Collector {
	completionTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sheepdog_runner",
		Name:      "request_duration_seconds",
		Help:      "Duration of the request.",
		Buckets:   []float64{0.1, 0.2, 0.3},
	}, []string{"method", "status", "location"})

	completionTime.With(prometheus.Labels{
		"method":   r.httpMethod,
		"status":   r.statusCode,
		"location": r.location.String(),
	}).Observe(r.latency)
	return completionTime
}

func (r *resultsHttp) GetCloudWatchDimensions() []cloudwatchtype.Dimension {
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

func (r *resultsHttp) SetError(err error) {
	r.errorProbe = err
}

func (r *resultsHttp) SetLatency(latency int64) {
	slog.Debug("Setting latency", "latency", latency)
	r.latency = float64(latency) / 1000.0
}

func (r *resultsHttp) SetStatusCode(statusCode string) {
	r.statusCode = statusCode
}

func (r *resultsHttp) GetErrorProbe() error {
	return r.errorProbe
}

func (r *resultsHttp) SetTime(time time.Time) {
	r.time = time
}

func (r *resultsHttp) MarshalJSON() ([]byte, error) {
	result := &ResultsHttpJSON{
		Id:         r.id,
		Time:       r.time.Format(time.RFC3339),
		Location:   r.location.String(),
		Latency:    fmt.Sprintf("%f", r.latency),
		ErrorProbe: fmt.Sprintf("%v", r.errorProbe),
		HttpMethod: r.httpMethod,
		StatusCode: r.statusCode,
	}
	return json.Marshal(result)
}
