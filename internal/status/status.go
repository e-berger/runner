package status

import "time"

type Status struct {
	Time    time.Time `json:"time"`
	ProbeId string    `json:"probe_id"`
	Status  State     `json:"status"`
	Details string    `json:"details"`
}

func NewStatus(started time.Time, probeId string, status State, details string) *Status {
	return &Status{
		Time:    started,
		ProbeId: probeId,
		Status:  status,
		Details: details,
	}
}
