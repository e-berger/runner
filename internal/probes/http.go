package probes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

type httpProbe struct {
	Probe
	httpProbeData
}

type httpProbeData struct {
	HttpMethod string `json:"method"`
	Url        string `json:"url"`
}

func NewHttpProbe(p *Probe) (IProbe, error) {
	var d httpProbeData
	err := json.Unmarshal([]byte(p.Data), &d)
	if err != nil {
		slog.Error("Error unmarshalling httptest data", "data", p.Data)
		return nil, err
	}

	h := &httpProbe{
		Probe: Probe{
			Id:       p.Id,
			Type:     HTTP,
			Location: p.Location,
		},
		httpProbeData: httpProbeData{
			HttpMethod: d.HttpMethod,
			Url:        d.Url,
		},
	}
	return h, nil
}

func (t *httpProbe) Launch() (metrics.IMetrics, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	req, err := http.NewRequest(t.HttpMethod, t.Url, nil)
	if err != nil {
		return nil, err
	}

	time_start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result := metrics.NewResultHttpDetails(
		t.Id,
		int(t.Location),
		time.Since(time_start).Milliseconds(),
		1,
		t.HttpMethod,
		resp.StatusCode)

	return result, nil
}

func (t *httpProbe) String() string {
	return fmt.Sprintf("http probe %s", t.Id)
}

func (t *httpProbe) GetType() ProbeType {
	return t.Type
}
