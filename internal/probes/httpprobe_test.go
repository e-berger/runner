package probes

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/stretchr/testify/assert"
)

func TestHttpProbe_Launch(t *testing.T) {

	probeHttpInfo := probes.NewHttpProbeInfo(http.MethodGet,
		"https://example.com")

	probe, err := probes.NewProbe("test", types.ThirtySeconds, types.Locations{types.EuropeLocation}, false, probeHttpInfo)
	assert.NoError(t, err)

	httprobe, err := NewHttpProbe(probe, types.EuropeLocation)
	assert.NoError(t, err)

	client := HTTPClientMock{}
	client.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil
	}
	httprobe.Launch(client)
	// assert.Equal(t, "test", result.GetId())
}
