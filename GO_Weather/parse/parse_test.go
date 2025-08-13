package parse

import "testing"

func TestParseUnits(t *testing.T) {
	testCases := []struct {
		unit    string
		expect  string
		wantErr bool
	}{
		{"CELSIUS", "metric", false},
		{"CELsIUS", "metric", false},
		{"celsius", "metric", false},
		{"c", "metric", false},
		{"C", "metric", false},
		{"metric", "metric", false},
		{"metriC", "metric", false},
		{"FAHRENHEIT", "imperial", false},
		{"FAHRENHEIt", "imperial", false},
		{"fahrenheit", "imperial", false},
		{"f", "imperial", false},
		{"F", "imperial", false},
		{"imperial", "imperial", false},
		{"IMPerial", "imperial", false},
		{"KELVIN", "standard", false},
		{"kelvin", "standard", false},
		{"KElvIN", "standard", false},
		{"k", "standard", false},
		{"K", "standard", false},
		{"standard", "standard", false},
		{"STandard", "standard", false},
		{"standart", "", true},
		{"", "", true},
		{"kesawa", "", true},
	}

	for _, testCase := range testCases {
		res, err := ParseUnits(testCase.unit)
		if (err != nil) != testCase.wantErr {
			t.Errorf("unexpected error for %v: %v", testCase, err)
		}
		if res != testCase.expect {
			t.Errorf("ParseUnits(%v) = %q; want %q", testCase.unit, res, testCase.expect)
		}
	}
}

func TestGetChar(t *testing.T) {
	testCases := []struct {
		unit    string
		expect  rune
		wantErr bool
	}{
		{"metric", 'C', false},
		{"metriC", '?', true},
		{"metris", '?', true},
		{"imperial", 'F', false},
		{"IMPerial", '?', true},
		{"impervial", '?', true},
		{"standard", 'K', false},
		{"STandard", '?', true},
		{"standart", '?', true},
		{"", '?', true},
		{"kesawa", '?', true},
	}

	for _, testCase := range testCases {
		res, err := GetChar(testCase.unit)
		if (err != nil) != testCase.wantErr {
			t.Errorf("unexpected error for %v: %v", testCase, err)
		}
		if res != testCase.expect {
			t.Errorf("GetChar(%v) = %q; want %q", testCase.unit, res, testCase.expect)
		}
	}
}
