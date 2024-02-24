package monitoring

import "fmt"

type MethodType uint

const (
	UNKNOWNMETHOD = iota
	PING
	HTTP
	HEAD
	TCP
	DNS
	SMTP
	SSH
)

const (
	PINGMethodString = "ping"
	HTTPMethodString = "http"
	HEADMethodString = "head"
	TCPMethodString  = "tcp"
	DNSMethodString  = "dns"
	SMTPMethodString = "smtp"
	SSHMethodString  = "ssh"
)

func (m MethodType) String() string {
	switch m {
	case PING:
		return PINGMethodString
	case HTTP:
		return HTTPMethodString
	case HEAD:
		return HEADMethodString
	case TCP:
		return TCPMethodString
	case DNS:
		return DNSMethodString
	case SMTP:
		return SMTPMethodString
	case SSH:
		return SSHMethodString
	default:
		panic("unhandled default case")
	}
}

// ParseLocation parses a method string into a MethodType
func ParseMethod(method string) (MethodType, error) {
	switch method {
	case PINGMethodString:
		return PING, nil
	case HTTPMethodString:
		return HTTP, nil
	case HEADMethodString:
		return HEAD, nil
	case TCPMethodString:
		return TCP, nil
	case DNSMethodString:
		return DNS, nil
	case SMTPMethodString:
		return SMTP, nil
	case SSHMethodString:
		return SMTP, nil
	}
	return UNKNOWNMETHOD, fmt.Errorf("unknown method: %s", method)
}
