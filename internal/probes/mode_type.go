package probes

import "fmt"

type Mode uint

const (
	UNKNOWNMODE Mode = iota
	CRON
	DIRECT
)

const (
	CronString   = "cron"
	DirectString = "direct"
)

func (m Mode) String() string {
	switch m {
	case CRON:
		return CronString
	case DIRECT:
		return DirectString
	default:
		panic("unhandled default case")
	}
}

// ParseMode parses a mode string into a Mode
func ParseMode(mode string) (Mode, error) {
	switch mode {
	case CronString:
		return CRON, nil
	case DirectString:
		return DIRECT, nil
	case "":
		return CRON, nil
	}
	return UNKNOWNMODE, fmt.Errorf("unknown mode: %s", mode)
}

func (m *Mode) UnmarshalJSON(data []byte) error {
	mode, err := ParseMode(string(data))
	if err != nil {
		return err
	}
	*m = mode
	return nil
}
