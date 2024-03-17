package probes

import "fmt"

type State uint

const (
	UNKNOWNSTATE State = iota
	UP
	ERROR
)

const (
	UpString    = "up"
	ErrorString = "error"
)

func (s State) String() string {
	switch s {
	case UP:
		return UpString
	case ERROR:
		return ErrorString
	default:
		panic("unhandled default case")
	}
}

func ParseState(state string) (State, error) {
	switch state {
	case UpString:
		return UP, nil
	case ErrorString:
		return ERROR, nil
	case "":
		return UP, nil
	default:
		return UNKNOWNSTATE, fmt.Errorf("unknown state: %s", state)
	}
}
