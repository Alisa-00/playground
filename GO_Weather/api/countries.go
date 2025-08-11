package api

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/countries.json
var countriesJSON []byte
var CountryCodes map[string]string

func init() {
	_ = json.Unmarshal(countriesJSON, &CountryCodes)
}
