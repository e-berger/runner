package status

type Status struct {
	ProbeId string `json:"probe_id"`
	Status  string `json:"status"`
	Error   bool   `json:"error"`
}

func NewStatus(probeId string, status string, error bool) *Status {
	return &Status{
		ProbeId: probeId,
		Status:  status,
		Error:   error,
	}
}
