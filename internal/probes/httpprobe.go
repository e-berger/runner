package probes

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

type httpProbe struct {
	domain.Probe
	Location types.Location
}

func NewHttpProbe(probe domain.Probe, location types.Location) (IProbe, error) {
	return httpProbe{
		probe,
		location,
	}, nil
}

func (t httpProbe) String() string {
	return fmt.Sprintf("http probe %s", t.Probe.GetId())
}

func (t httpProbe) GetId() string {
	return t.Probe.GetId()
}

func (t httpProbe) IsInError() bool {
	return t.Probe.IsInError()
}

func (t httpProbe) Launch() (metrics.IMetrics, error) {
	// Redirect
	var checkRedirect func(req *http.Request, via []*http.Request) error
	if t.GetHttpProbeInfo().FollowRedirects {
		checkRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}
	client := &http.Client{
		CheckRedirect: checkRedirect,
	}

	// Timeout
	timeout := time.Duration(t.Probe.GetHttpProbeInfo().Timeout)
	slog.Debug("Probe timeout", "second", timeout.Seconds())

	req, err := http.NewRequest(t.GetHttpProbeInfo().Method, t.GetHttpProbeInfo().Url, nil)
	if err != nil {
		return nil, err
	}

	// Headers
	for k, v := range t.GetHttpProbeInfo().Headers {
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
		t.GetId(),
		int(t.Location),
		time_start,
		time.Since(time_start).Milliseconds(),
		DEFAULT_VALID,
		t.GetHttpProbeInfo().Method,
		resp.StatusCode)

	// Analyse the response with constraints
	return result, t.analyse(resp)
}

func (t httpProbe) analyse(resp *http.Response) error {
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

func (t httpProbe) validExpectedHeaders(headers http.Header) error {
	slog.Info("Probe headers", "probe", t.GetId(), "headers", headers, "expected", t.GetHttpProbeInfo().ExpectedHeaders)
	for k, v := range t.GetHttpProbeInfo().ExpectedHeaders {
		if headers.Get(k) != v {
			return fmt.Errorf("unexpected header %s: %s", k, v)
		}
	}
	return nil
}

func (t httpProbe) validExpectedContent(body io.ReadCloser) error {
	if len(t.GetHttpProbeInfo().ExpectedContent) > 0 {
		b, _ := io.ReadAll(body)
		content := string(b)
		slog.Info("Probe content", "probe", t.GetId(), "content", content, "expected", t.GetHttpProbeInfo().ExpectedContent)
		match, err := regexp.MatchString(t.GetHttpProbeInfo().ExpectedContent, content)
		if err != nil || !match {
			return fmt.Errorf("unexpected content %s", t.GetHttpProbeInfo().ExpectedContent)
		}
	}
	return nil
}

func (t httpProbe) validExpectedStatus(statusCode int) error {
	slog.Info("Probe status code", "probe", t.GetId(), "status", statusCode, "expected", t.GetHttpProbeInfo().ExpectedStatusCode, "not_expected", t.GetHttpProbeInfo().NotExpectedStatusCode)
	if len(t.GetHttpProbeInfo().ExpectedStatusCode) > 0 && !matchStatus(t.GetHttpProbeInfo().ExpectedStatusCode, statusCode) {
		return fmt.Errorf("unexpected status code %d", statusCode)
	} else if len(t.GetHttpProbeInfo().NotExpectedStatusCode) > 0 && matchStatus(t.GetHttpProbeInfo().NotExpectedStatusCode, statusCode) {
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
