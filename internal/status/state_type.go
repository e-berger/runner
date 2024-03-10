package status

import "fmt"

type State uint

const (
	UNKNOWNLOCATION State = iota
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
	default:
		return UNKNOWNLOCATION, fmt.Errorf("unknown state: %s", state)
	}
}
