package api

import "testing"

func TestCountryCodes(t *testing.T) {
	testCases := []struct {
		country string
		expect  string
	}{
		{"france", "FR"},
		{"bangladesh", "BD"},
		{"belgium", "BE"},
		{"burkina faso", "BF"},
		{"bulgaria", "BG"},
		{"bosnia and herzegovina", "BA"},
		{"barbados", "BB"},
		{"wallis and futuna", "WF"},
		{"saint barthelemy", "BL"},
		{"ukraine", "UA"},
		{"qatar", "QA"},
		{"mozambique", "MZ"},
		{"QATAR", ""},
		{"kewa", ""},
	}

	for _, testCase := range testCases {

		code := CountryCodes[testCase.country]
		if code != testCase.expect {
			t.Errorf("CountryCodes[%v] = %v; want %q", testCase.country, code, testCase.expect)
		}
	}
}
