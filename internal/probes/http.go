package probes

import (
	"context"
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
	HttpMethod            string            `json:"method"`
	Url                   string            `json:"url"`
	Timeout               int               `json:"timeout"`
	Content               string            `json:"content"`
	Headers               map[string]string `json:"headers"`
	ExpectedStatusCode    []int             `json:"expected_status_code"`
	NotExpectedStatusCode []int             `json:"not_expected_status_code"`
	ExpectedContent       string            `json:"expected_content"`
	ExpectedHeaders       map[string]string `json:"expected_headers"`
	FollowRedirects       int               `json:"follow_redirects"`
}

func NewHttpProbe(p *Probe) (IProbe, error) {
	var d httpProbeData
	err := json.Unmarshal(p.Data, &d)
	if err != nil {
		slog.Error("Error unmarshalling httptest data", "data", p.Data)
		return nil, err
	}

	h := &httpProbe{
		Probe: Probe{
			Id:       p.Id,
			Type:     HTTP,
			Location: p.Location,
			Error:    p.Error,
		},
		httpProbeData: httpProbeData{
			HttpMethod:      d.HttpMethod,
			Url:             d.Url,
			Timeout:         d.Timeout,
			FollowRedirects: d.FollowRedirects,
		},
	}
	return h, nil
}

func (t *httpProbe) Launch() (metrics.IMetrics, error) {

	//redirect
	var checkRedirect func(req *http.Request, via []*http.Request) error
	if t.FollowRedirects > 0 {
		checkRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}
	client := &http.Client{
		CheckRedirect: checkRedirect,
	}

	//timeout
	timeout := time.Duration(t.Timeout)
	slog.Debug("Probe timeout", "second", timeout.Seconds())

	req, err := http.NewRequest(t.HttpMethod, t.Url, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.TODO())
	req = req.WithContext(ctx)

	time.AfterFunc(timeout, func() {
		cancel()
	})

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

func (p *httpProbe) GetId() string {
	return p.Id
}

func (p *httpProbe) GetError() bool {
	return p.Error
}
