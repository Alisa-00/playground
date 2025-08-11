package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/countries.json
var countriesJSON []byte
var CountryCodes map[string]string

func init() {
	err := json.Unmarshal(countriesJSON, &CountryCodes)
	if err != nil {
		panic(fmt.Errorf("failed to load countries.json: %w", err))
	}
}
