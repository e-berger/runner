package metrics

import (
	"fmt"
	"strconv"
)

type IMetrics interface {
	GetId() string
	GetLocation() string
	GetLatency() float64
	String() string
}

type Metrics struct {
	Id       string  `json:"id"`
	Location string  `json:"location"`
	Latency  float64 `json:"latency"`
	Valid    string  `json:"valid"`
}

func NewMetrics(id string, location int, latency int64, valid int) IMetrics {
	return &Metrics{
		Id:       id,
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
