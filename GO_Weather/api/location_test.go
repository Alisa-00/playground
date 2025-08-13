package api

import "testing"

func TestCityCountryCacheKey(t *testing.T) {
	testCases := []struct {
		CC      CityCountry
		expect  string
		wantErr bool
	}{
		{CityCountry{"Paris", "France"}, "paris,FR", false},
		{CityCountry{"PariS", "France"}, "paris,FR", false},
		{CityCountry{"Paris", "france"}, "paris,FR", false},
		{CityCountry{"Pari", "France"}, "pari,FR", false},
		{CityCountry{"Paris", "Franc"}, "paris", false},
		{CityCountry{"Paris", "united states"}, "paris,US", false},
		{CityCountry{"", ""}, "", true},
		{CityCountry{"", "UNITED states"}, "US", false},
	}

	for _, testCase := range testCases {
		res, err := testCase.CC.CacheKey()
		if (err != nil) != testCase.wantErr {
			t.Errorf("unexpected error for %v: %v", testCase, err)
		}
		if res != testCase.expect {
			t.Errorf("CacheKey(%v,%v) = %q; want %q", testCase.CC.City, testCase.CC.Country, res, testCase.expect)
		}
	}
}

func TestCityCountryQueryParam(t *testing.T) {
	testCases := []struct {
		CC      CityCountry
		expect  string
		wantErr bool
	}{
		{CityCountry{"Paris", "France"}, "q=Paris,France", false},
		{CityCountry{"PariS", "France"}, "q=PariS,France", false},
		{CityCountry{"Paris", "france"}, "q=Paris,france", false},
		{CityCountry{"Pari", "France"}, "q=Pari,France", false},
		{CityCountry{"Paris", "Franc"}, "q=Paris,Franc", false},
		{CityCountry{"Paris", "united states"}, "q=Paris,united+states", false},
		{CityCountry{"Paris", ""}, "q=Paris", false},
		{CityCountry{"", ""}, "", true},
		{CityCountry{"", "UNITED states"}, "q=UNITED+states", false},
	}

	for _, testCase := range testCases {
		res, err := testCase.CC.QueryParam()
		if (err != nil) != testCase.wantErr {
			t.Errorf("unexpected error for %v: %v", testCase, err)
		}
		if res != testCase.expect {
			t.Errorf("QueryParam(%v,%v) = %q; want %q", testCase.CC.City, testCase.CC.Country, res, testCase.expect)
		}
	}
}

func TestLatLonCacheKey(t *testing.T) {
	testCases := []struct {
		LL      LatLon
		expect  string
		wantErr bool
	}{
		{LatLon{0, 0}, "", true},
		{LatLon{1, 0}, "1.00-0.00", false},
		{LatLon{1, 1}, "1.00-1.00", false},
		{LatLon{1.000000000000001, 0}, "1.00-0.00", false},
		{LatLon{1.001, 0}, "1.00-0.00", false},
		{LatLon{1.1, 0}, "1.10-0.00", false},
		{LatLon{-1, -1}, "-1.00--1.00", false},
	}

	for _, testCase := range testCases {
		res, err := testCase.LL.CacheKey()
		if (err != nil) != testCase.wantErr {
			t.Errorf("unexpected error for %v: %v", testCase, err)
		}
		if res != testCase.expect {
			t.Errorf("CacheKey(%v,%v) = %q; want %q", testCase.LL.Lat, testCase.LL.Lon, res, testCase.expect)
		}
	}
}

func TestLatLonQueryParam(t *testing.T) {
	testCases := []struct {
		LL      LatLon
		expect  string
		wantErr bool
	}{
		{LatLon{0, 0}, "", true},
		{LatLon{1, 0}, "lat=1.000000&lon=0.000000", false},
		{LatLon{1, 1}, "lat=1.000000&lon=1.000000", false},
		{LatLon{1.000000000000001, 0}, "lat=1.000000&lon=0.000000", false},
		{LatLon{1.001, 0}, "lat=1.001000&lon=0.000000", false},
		{LatLon{1.1, 0}, "lat=1.100000&lon=0.000000", false},
		{LatLon{-1, -1}, "lat=-1.000000&lon=-1.000000", false},
	}

	for _, testCase := range testCases {
		res, err := testCase.LL.QueryParam()
		if (err != nil) != testCase.wantErr {
			t.Errorf("unexpected error for %v: %v", testCase, err)
		}
		if res != testCase.expect {
			t.Errorf("QueryParam(%v,%v) = %q; want %q", testCase.LL.Lat, testCase.LL.Lon, res, testCase.expect)
		}
	}
}
