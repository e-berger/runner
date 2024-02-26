package domain

import "fmt"

type Location uint

const (
	UNKNOWNLOCATION Location = iota
	NORTHAMERICA
	SOUTHAMERICA
	EUROPE
	ASIA
	AUSTRALIA
	AFRICA
)

const (
	// NorthAmericaLocationString is the string representation of NorthAmericaLocation
	NorthAmericaLocationString = "north_america"
	// SouthAmericaLocationString is the string representation of SouthAmericaLocation
	SouthAmericaLocationString = "south_america"
	// EuropeLocationString is the string representation of EuropeLocation
	EuropeLocationString = "europe"
	// AsiaLocationString is the string representation of AsiaLocation
	AsiaLocationString = "asia"
	// AustraliaLocationString is the string representation of AustraliaLocation
	AustraliaLocationString = "australia"
	// AfricaLocationString is the string representation of AfricaLocation
	AfricaLocationString = "africa"
)

func (l Location) String() string {
	switch l {
	case NORTHAMERICA:
		return NorthAmericaLocationString
	case SOUTHAMERICA:
		return SouthAmericaLocationString
	case EUROPE:
		return EuropeLocationString
	case ASIA:
		return AsiaLocationString
	case AUSTRALIA:
		return AustraliaLocationString
	case AFRICA:
		return AfricaLocationString
	default:
		panic("unhandled default case")
	}
}

// ParseLocation parses a location string into a Location
func ParseLocation(location string) (Location, error) {
	switch location {
	case NorthAmericaLocationString:
		return NORTHAMERICA, nil
	case SouthAmericaLocationString:
		return SOUTHAMERICA, nil
	case EuropeLocationString:
		return EUROPE, nil
	case AsiaLocationString:
		return ASIA, nil
	case AustraliaLocationString:
		return AUSTRALIA, nil
	case AfricaLocationString:
		return AFRICA, nil
	}
	return UNKNOWNLOCATION, fmt.Errorf("unknown location: %s", location)
}
