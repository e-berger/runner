package monitoring

import "fmt"

type MethodType uint

const (
	PING = iota
	HTTP
	HEAD
	TCP
	DNS
	SMTP
	SSH
)

var methodTypeStrings = []string{
	"ping",
	"http",
	"head",
	"tcp",
	"dns",
	"smtp",
	"ssh",
}

func (mt MethodType) name() string {
	return methodTypeStrings[mt]
}

func (mt MethodType) ordinal() int {
	return int(mt)
}

func (mt MethodType) values() *[]string {
	return &methodTypeStrings
}

func ValueOf(name string) (MethodType, error) {
	for i, n := range methodTypeStrings {
		if n == name {
			return MethodType(i), nil
		}
	}
	return 0, fmt.Errorf("MethodType %s not found", name)
}

func GetMethodType(i int) (MethodType, error) {
	if i < 0 || i >= len(methodTypeStrings) {
		return 0, fmt.Errorf("MethodType %d not found", i)
	}
	return MethodType(i), nil
}

func (mt MethodType) String() string {
	return methodTypeStrings[mt]
}
