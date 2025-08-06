package parse

import (
	"fmt"
	"strings"
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
