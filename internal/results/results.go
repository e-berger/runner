package results

import (
	"time"

	cloudwatchtype "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/prometheus/client_golang/prometheus"
)

type IResults interface {
	GetCloudWatchDimensions() []cloudwatchtype.Dimension
	GetPrometheusMetrics() prometheus.Collector
	String() string
	GetId() string
	GetLatency() float64
	GetTime() time.Time
	SetError(err error)
	GetErrorProbe() error
	MarshalJSON() ([]byte, error)
}

type results struct {
	id         string
	time       time.Time
	location   types.Location
	latency    float64
	code       types.Code
	errorProbe error
}
