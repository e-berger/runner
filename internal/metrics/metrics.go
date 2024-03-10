package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type IMetrics interface {
	GetId() string
	GetLocation() string
	GetLatency() float64
	GetTime() time.Time
	String() string
	GetMetrics() prometheus.Collector
}

type Metrics struct {
	Id       string    `json:"id"`
	Time     time.Time `json:"time"`
	Location string    `json:"location"`
	Latency  float64   `json:"latency"`
	Valid    string    `json:"valid"`
}

func NewMetrics(id string, location int, latency int64, valid int) IMetrics {
	return &Metrics{
		Id:       id,
		Time:     time.Now(),
		Location: strconv.Itoa(location),
		Latency:  float64(latency) / 1000.0,
		Valid:    strconv.Itoa(valid),
	}
}

func (r *Metrics) GetId() string {
	return r.Id
}

func (r *Metrics) GetLocation() string {
	return r.Location
}

func (r *Metrics) GetLatency() float64 {
	return r.Latency
}

func (r *Metrics) String() string {
	return fmt.Sprintf("Id: %s, Location: %s, Latency: %f, Valid: %s", r.Id, r.Location, r.Latency, r.Valid)
}

func (r *Metrics) GetMetrics() prometheus.Collector {
	return nil
}

func (r *Metrics) GetTime() time.Time {
	return r.Time
}
