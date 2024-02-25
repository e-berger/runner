package probes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/e-berger/sheepdog-runner/src/internal/domain"
	"github.com/e-berger/sheepdog-runner/src/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type httpProbe struct {
	domain.Probe
	httpProbeData
	result *metrics.MetricsHttp
}

type httpProbeData struct {
	HttpMethod string `json:"method"`
	Url        string `json:"url"`
}

func NewHttpProbe(p *domain.Probe) (domain.IProbe, error) {
	var d httpProbeData
	err := json.Unmarshal([]byte(p.Data), &d)
	if err != nil {
		slog.Error("Error unmarshalling httptest data", "data", p.Data)
		return nil, err
	}

	h := &httpProbe{
		Probe: domain.Probe{
			Id:       p.Id,
			Type:     domain.HTTP,
			Location: p.Location,
		},
		httpProbeData: httpProbeData{
			HttpMethod: d.HttpMethod,
			Url:        d.Url,
		},
	}
	return h, nil
}

func (t *httpProbe) Launch() error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	req, err := http.NewRequest(t.HttpMethod, t.Url, nil)
	if err != nil {
		return err
	}

	time_start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	t.result = metrics.NewResultHttpDetails(
		t.Id,
		int(t.Location),
		time.Since(time_start).Milliseconds(),
		1,
		t.HttpMethod,
		resp.StatusCode)

	return nil
}

func (t *httpProbe) Push(pushGateway *metrics.Push) error {
	completionTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sheepdog_runner",
		Name:      "request_duration_seconds",
		Help:      "Duration of the request.",
		Buckets:   []float64{0.1, 0.2, 0.3},
	}, []string{"method", "status"})

	completionTime.With(prometheus.Labels{"method": t.HttpMethod, "status": t.result.GeStatusCode()}).Observe(t.result.GetLatency())
	return pushGateway.Send(fmt.Sprint(t.Id), completionTime)
}

func (t *httpProbe) String() string {
	return fmt.Sprintf("http test %d", t.Id)
}

func (t *httpProbe) GetType() domain.ProbeType {
	return t.Type
}

func (t *httpProbe) GetResult() metrics.IMetrics {
	return t.result
}
