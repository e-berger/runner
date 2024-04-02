package status

import (
	"time"

	"github.com/e-berger/sheepdog-runner/internal/probes"
)

type Status struct {
	Time     time.Time       `json:"time"`
	ProbeId  string          `json:"probe_id"`
	State    probes.State    `json:"status"`
	Details  string          `json:"details"`
	Mode     probes.Mode     `json:"mode"`
	Location probes.Location `json:"location"`
}

func NewStatus(started time.Time, probeId string, state probes.State, details string, mode probes.Mode, location probes.Location) *Status {
	if mode == probes.UNKNOWNMODE {
		mode = probes.CRON
	}
	return &Status{
		Time:     started,
		ProbeId:  probeId,
		State:    state,
		Details:  details,
		Mode:     mode,
		Location: location,
	}
}
