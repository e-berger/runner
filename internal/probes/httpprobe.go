package probes

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/results"
)

const default_timeout = 10

type httpProbe struct {
	domain.Probe
	location types.Location
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

func (t httpProbe) GetHttpClient() HTTPClient {
	// Redirect
	var checkRedirect func(req *http.Request, via []*http.Request) error
	if t.GetHttpProbeInfo().FollowRedirects {
		checkRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}

	// Default timeout : 10s
	timeout := time.Duration(default_timeout * time.Second)
	if t.Probe.GetHttpProbeInfo().Timeout != 0 {
		timeout = time.Duration(t.Probe.GetHttpProbeInfo().Timeout)
		// Timeout can't be lower than 1s
		if timeout < (1 * time.Second) {
			timeout = 1 * time.Second
		}
	}

	// Allow insecure & timeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: t.GetHttpProbeInfo().AllowInsecure},
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}

	return &http.Client{
		CheckRedirect: checkRedirect,
		Transport:     tr,
		Timeout:       timeout,
	}
}

func (t httpProbe) Launch(client HTTPClient) results.IResults {
	result := results.NewResultsHttpEmpty(
		t.GetId(),
		t.location,
		t.GetHttpProbeInfo().Method)

	// Body
	var bodyReader io.Reader
	if t.GetHttpProbeInfo().Content != "" {
		bodyReader = bytes.NewReader([]byte(t.GetHttpProbeInfo().Content))
	}

	req, err := http.NewRequest(t.GetHttpProbeInfo().Method, t.GetHttpProbeInfo().Url, bodyReader)
	if err != nil {
		result.SetCode(types.ERROR)
		result.SetError(err)
		return result
	}

	// Headers
	for k, v := range t.GetHttpProbeInfo().Headers {
		req.Header.Add(k, v)
	}

	req = req.WithContext(context.Background())

	time_start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			result.SetCode(types.TIMEOUT)
		} else {
			result.SetCode(types.ERROR)
		}
		result.SetError(err)
		return result
	}
	defer resp.Body.Close()

	result.SetTime(time_start)
	result.SetCode(types.UP)
	result.SetLatency(time.Since(time_start).Milliseconds())
	// Analyse the response with constraints
	errAnalyse := t.analyse(resp)
	if errAnalyse != nil {
		result.SetCode(types.ERROR)
	}
	result.SetError(errAnalyse)

	return result
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
	slog.Info("Probe status code", "probe", t.GetId(), "status", statusCode, "expected", t.GetHttpProbeInfo().ExpectedStatusCodes, "not_expected", t.GetHttpProbeInfo().NotExpectedStatusCodes)
	if len(t.GetHttpProbeInfo().ExpectedStatusCodes) > 0 {
		if t.GetHttpProbeInfo().ExpectedStatusCodes.IsValid() {
			if !matchStatus(t.GetHttpProbeInfo().ExpectedStatusCodes, statusCode) {
				return fmt.Errorf("unexpected status code %d", statusCode)
			}
		} else {
			slog.Error("Invalid status code family", "probe", t.GetId(), "status", t.GetHttpProbeInfo().ExpectedStatusCodes)
		}
	} else if len(t.GetHttpProbeInfo().NotExpectedStatusCodes) > 0 {
		if t.GetHttpProbeInfo().NotExpectedStatusCodes.IsValid() {
			if matchStatus(t.GetHttpProbeInfo().NotExpectedStatusCodes, statusCode) {
				return fmt.Errorf("unexpected status code %d", statusCode)
			}
		} else {
			slog.Error("Invalid status code family", "probe", t.GetId(), "status", t.GetHttpProbeInfo().NotExpectedStatusCodes)
		}
	}
	return nil
}

func matchStatus(s types.HttpStatusCodeFamilies, code int) bool {
	for _, v := range s.ToHttpStatusCode() {
		if v == code {
			return true
		}
	}
	return false
}
