package parse

import (
	"fmt"
	"strings"
	"weather/api"
)

var chars = map[string]rune{
	"metric":   'C',
	"imperial": 'F',
	"standard": 'K',
}
var unitAliases = map[string]string{
	"metric":     "metric",
	"celsius":    "metric",
	"c":          "metric",
	"imperial":   "imperial",
	"fahrenheit": "imperial",
	"f":          "imperial",
	"standard":   "standard",
	"kelvin":     "standard",
	"k":          "standard",
}

func ParseUnits(units string) (string, error) {
	unit, ok := unitAliases[strings.ToLower(units)]
	if !ok {
		return unit, fmt.Errorf("invalid units")
	}
	return unit, nil
}

func GetChar(units string) (rune, error) {
	unitChar, ok := chars[units]

	if !ok {
		return '?', fmt.Errorf("no unit char for the unit: %s", units)
	}

	return unitChar, nil
}

func GetLocation(city string, country string, lat float64, lon float64) (api.Location, error) {

	if lat != 0 || lon != 0 {
		return api.Location{Latitude: lat, Longitude: lon}, nil
	}

	if city != "" {
		if country != "" {
			return api.Location{City: city, Country: country}, nil
		}
		return api.Location{City: city}, nil
	}

	if country != "" {
		return api.Location{Country: country}, nil
	}

	return api.Location{}, fmt.Errorf("invalid input. city, country or latitude and longitude have to be valid")
}
