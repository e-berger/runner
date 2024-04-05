package metrics

import (
	"time"

	cloudwatchtype "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/prometheus/client_golang/prometheus"
)

type IMetrics interface {
	GetCloudWatchDimensions() []cloudwatchtype.Dimension
	GetPrometheusMetrics() prometheus.Collector
	String() string
	GetId() string
	GetLatency() float64
	GetTime() time.Time
}

type metrics struct {
	id       string
	time     time.Time
	location types.Location
	latency  float64
	valid    string
}
