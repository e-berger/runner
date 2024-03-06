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

// "location": "australia",
// "items": [
//     {
//         "id": "2d2f35Ry74DX9F9piVm4FWUuz3b",
//         "info": {
//             "timeout": 10000000000,
//             "method": "GET",
//             "url": "https://observations-service-api.eu.finalcad.cloud/healthz/live",
//             "expected_status_code": [
//                 200
//             ]
//         },
//         "type": 2
//     }
// ]

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
	FollowRedirect        int               `json:"follow_redirect"`
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
		},
		httpProbeData: httpProbeData{
			HttpMethod:     d.HttpMethod,
			Url:            d.Url,
			Timeout:        d.Timeout,
			FollowRedirect: d.FollowRedirect,
		},
	}
	return h, nil
}

func (t *httpProbe) Launch() (metrics.IMetrics, error) {

	//Redirect
	var checkRedirect func(req *http.Request, via []*http.Request) error
	if t.FollowRedirect > 0 {
		checkRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}
	client := &http.Client{
		CheckRedirect: checkRedirect,
	}

	//timeout
	time_start := time.Now()
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(time.Duration(t.Timeout)*time.Second, func() {
		cancel()
	})

	req, err := http.NewRequest(t.HttpMethod, t.Url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

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
