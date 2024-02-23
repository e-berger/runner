package monitor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/e-berger/sheepdog-runner/src/internal/monitoring"
)

type httpTest struct {
	Monitor
	httpTestData
}

type httpTestData struct {
	HttpMethod string `json:"method"`
	Url        string `json:"url"`
}

func newHttpTest(m InfosMonitor) (IMonitor, error) {
	var d httpTestData
	err := json.Unmarshal([]byte(m.Data), &d)
	if err != nil {
		slog.Error("Error unmarshalling httptest data", "data", m.Data)
		return nil, err
	}

	h := &httpTest{
		Monitor: Monitor{
			id:     m.Id,
			method: monitoring.HTTP,
		},
		httpTestData: httpTestData{
			HttpMethod: d.HttpMethod,
			Url:        d.Url,
		},
	}

	return h, nil
}

func (t *httpTest) GetMethod() monitoring.MethodType {
	return t.method
}

func (t *httpTest) Launch() error {
	time_start := time.Now()
	resp, err := http.Get(t.Url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		slog.Error("Error on http ping", "error", err, "httptest", t)
		return err
	}
	fmt.Println(time.Since(time_start), t.Url)
	return nil
}

func (t *httpTest) String() string {
	return fmt.Sprintf("http test %d", t.id)
}
