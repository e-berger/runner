package results

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudwatchtype "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/prometheus/client_golang/prometheus"
)

type resultsHttp struct {
	results
	httpMethod string
}

type ResultsHttpJSON struct {
	Id         string `json:"probe_id"`
	Time       string `json:"time"`
	Location   string `json:"location"`
	Latency    string `json:"latency"`
	ErrorProbe string `json:"error"`
	Code       string `json:"code"`
	HttpMethod string `json:"httpMethod"`
}

func (r *resultsHttp) String() string {
	return fmt.Sprintf("id: %s, location: %s, latency: %f, httpMethod: %s, Code: %s",
		r.id, r.location, r.latency, r.httpMethod, r.code.String())
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
			code:     types.Code(statusCode),
		},
		httpMethod: httpMethod,
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
	}
}

func (r *resultsHttp) SetError(err error) {
	r.errorProbe = err
}

func (r *resultsHttp) SetLatency(latency int64) {
	slog.Debug("Result latency", "latency", latency, "probe", r.id)
	r.latency = float64(latency) / 1000.0
}

func (r *resultsHttp) SetCode(code types.Code) {
	r.code = code
}

func (r *resultsHttp) GetErrorProbe() error {
	return r.errorProbe
}

func (r *resultsHttp) SetTime(time time.Time) {
	r.time = time
}

func (r *resultsHttp) GetCode() types.Code {
	return r.code
}

func (r *resultsHttp) MarshalJSON() ([]byte, error) {
	result := &ResultsHttpJSON{
		Id:         r.id,
		Time:       r.time.Format(time.RFC3339),
		Location:   r.location.String(),
		Latency:    fmt.Sprintf("%f", r.latency),
		ErrorProbe: fmt.Sprintf("%v", r.errorProbe),
		HttpMethod: r.httpMethod,
		Code:       r.code.String(),
	}
	return json.Marshal(result)
}
