package metrics

import (
	"fmt"
	"strconv"
)

type MetricsHttp struct {
	Metrics
	HttpMethod string `json:"http_method"`
	StatusCode string `json:"status_code"`
}

func (r *MetricsHttp) GetId() string {
	return r.Id
}

func (r *MetricsHttp) GetLocation() string {
	return r.Location
}

func (r *MetricsHttp) GetLatency() float64 {
	return r.Latency
}

func (r *MetricsHttp) GeStatusCode() string {
	return r.StatusCode
}

func (r *MetricsHttp) String() string {
	return fmt.Sprintf("Id: %s, Location: %s, Latency: %f, Valid: %s, HttpMethod: %s, StatusCode: %s", r.Id, r.Location, r.Latency, r.Valid, r.HttpMethod, r.StatusCode)
}

func NewResultHttpDetails(id string, location int, latency int64, valid int, httpMethod string, statusCode int) *MetricsHttp {
	return &MetricsHttp{
		Metrics: Metrics{
			Id:       id,
			Location: strconv.Itoa(location),
			Latency:  float64(latency) / 1000.0,
			Valid:    strconv.Itoa(valid),
		},
		HttpMethod: httpMethod,
		StatusCode: strconv.Itoa(statusCode),
	}
}
