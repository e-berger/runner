package probes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
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
			State:    p.State,
		},
		httpProbeData: httpProbeData{
			HttpMethod:            d.HttpMethod,
			Url:                   d.Url,
			Timeout:               d.Timeout,
			Headers:               d.Headers,
			Content:               d.Content,
			FollowRedirects:       d.FollowRedirects,
			ExpectedStatusCode:    d.ExpectedStatusCode,
			NotExpectedStatusCode: d.NotExpectedStatusCode,
			ExpectedContent:       d.ExpectedContent,
			ExpectedHeaders:       d.ExpectedHeaders,
		},
	}
	return h, nil
}

func (t *httpProbe) Launch() (metrics.IMetrics, error) {

	// Redirect
	var checkRedirect func(req *http.Request, via []*http.Request) error
	if t.FollowRedirects > 0 {
		checkRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}
	client := &http.Client{
		CheckRedirect: checkRedirect,
	}

	// Timeout
	timeout := time.Duration(t.Timeout)
	slog.Debug("Probe timeout", "second", timeout.Seconds())

	req, err := http.NewRequest(t.HttpMethod, t.Url, nil)
	if err != nil {
		return nil, err
	}

	// Headers
	for k, v := range t.Headers {
		req.Header.Add(k, v)
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
	defer resp.Body.Close()

	result := metrics.NewResultHttpDetails(
		t.Id,
		int(t.Location),
		time_start,
		time.Since(time_start).Milliseconds(),
		1,
		t.HttpMethod,
		resp.StatusCode)

	// Analyse the response with constraints
	err = t.analyse(resp)

	return result, err
}

func (t *httpProbe) String() string {
	return fmt.Sprintf("http probe %s", t.Id)
}

func (t *httpProbe) GetType() ProbeType {
	return t.Type
}

func (t *httpProbe) GetId() string {
	return t.Id
}

func (t *httpProbe) IsError() bool {
	return t.State == ERROR
}

func (t *httpProbe) analyse(resp *http.Response) error {
	if err := t.validExpectedStatus(resp.StatusCode); err != nil {
		return err
	}
	if err := t.validExpectedContent(resp.Body); err != nil {
		return err
	}
	if err := t.validExpectedHeaders(resp.Header); err != nil {
		return err
	}

	return nil
}

func (t *httpProbe) validExpectedHeaders(headers http.Header) error {
	slog.Info("Probe headers", "probe", t.Id, "headers", headers, "expected", t.ExpectedHeaders)
	for k, v := range t.ExpectedHeaders {
		if headers.Get(k) != v {
			return fmt.Errorf("unexpected header %s: %s", k, v)
		}
	}
	return nil
}

func (t *httpProbe) validExpectedContent(body io.ReadCloser) error {
	if len(t.ExpectedContent) > 0 {
		b, _ := io.ReadAll(body)
		content := string(b)
		slog.Info("Probe content", "probe", t.Id, "content", content, "expected", t.ExpectedContent)
		match, err := regexp.MatchString(t.ExpectedContent, content)
		if err != nil || !match {
			return fmt.Errorf("unexpected content %s", t.ExpectedContent)
		}
	}
	return nil
}

func (t *httpProbe) validExpectedStatus(statusCode int) error {
	slog.Info("Probe status code", "probe", t.Id, "status", statusCode, "expected", t.ExpectedStatusCode, "not_expected", t.NotExpectedStatusCode)
	if len(t.ExpectedStatusCode) > 0 && !matchStatus(t.ExpectedStatusCode, statusCode) {
		return fmt.Errorf("unexpected status code %d", statusCode)
	} else if len(t.NotExpectedStatusCode) > 0 && matchStatus(t.NotExpectedStatusCode, statusCode) {
		return fmt.Errorf("unexpected status code %d", statusCode)
	}
	return nil
}

func matchStatus(s []int, code int) bool {
	for _, v := range s {
		if v == code {
			return true
		}
	}
	return false
}
