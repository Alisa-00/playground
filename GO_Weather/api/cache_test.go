package api

import (
	"testing"
	"time"
)

func TestPutRead(t *testing.T) {

	cache, err := LoadCacheFile()
	if err != nil {
		t.Fatalf("LoadCacheFile() failed: %v", err)
	}
	time1, err := time.Parse(time.RFC3339, "2025-08-12T18:25:10+03:00")
	if err != nil {
		t.Fatalf("time.Parse(time1) failed: %v", err)
	}
	time2, err := time.Parse(time.RFC3339, "2025-08-12T18:31:03+03:00")
	if err != nil {
		t.Fatalf("time.Parse(time2) failed: %v", err)
	}

	testCases := []struct {
		loc       Location
		weather   Weather
		expectErr bool
	}{
		{CityCountry{"Paris", "France"}, Weather{QueryType(Current), "Paris", "FR", 0, 0, []WeatherDay{{35.3, 34.49, "clear sky", time1}}}, false},
		{CityCountry{"Paris", "america"}, Weather{QueryType(Current), "Paris", "US", 0, 0, []WeatherDay{{28.12, 30.76, "clear sky", time2}}}, false},
		{LatLon{48.8534, 2.3488}, Weather{QueryType(Current), "Paris", "FR", 48.8534, 2.3488, []WeatherDay{{35.3, 34.49, "clear sky", time1}}}, false},
		{LatLon{33.6609, -95.5555}, Weather{QueryType(Current), "Paris", "US", 33.6609, -95.5555, []WeatherDay{{28.12, 30.76, "clear sky", time2}}}, false},
		{CityCountry{"Paris", "France"}, Weather{QueryType(Forecast), "Paris", "FR", 0, 0, []WeatherDay{{35.3, 34.49, "clear sky", time1}}}, false},
		{CityCountry{"Paris", "america"}, Weather{QueryType(Forecast), "Paris", "US", 0, 0, []WeatherDay{{28.12, 30.76, "clear sky", time2}}}, false},
		{LatLon{48.8534, 2.3488}, Weather{QueryType(Forecast), "Paris", "FR", 48.8534, 2.3488, []WeatherDay{{35.3, 34.49, "clear sky", time1}}}, false},
		{LatLon{33.6609, -95.5555}, Weather{QueryType(Forecast), "Paris", "US", 33.6609, -95.5555, []WeatherDay{{28.12, 30.76, "clear sky", time2}}}, false},
	}

	for _, testCase := range testCases {

		err = cache.Put(testCase.loc, testCase.weather)
		if (err != nil) != testCase.expectErr {
			t.Fatalf("error on put: %v", err)
		}

		queryType := "current"
		if testCase.weather.Type != QueryType(Current) {
			queryType = "forecast"
		}

		weth, err := cache.Read(testCase.loc, queryType)
		if (err != nil) != testCase.expectErr {
			t.Fatalf("error on read: %v", err)
		}

		switch tp := testCase.loc.(type) {
		case CityCountry:
			{
				if (weth.City != testCase.weather.City) || (weth.Country != testCase.weather.Country) || (weth.Type != testCase.weather.Type) || (len(weth.List) != len(testCase.weather.List)) {
					t.Errorf("Mismatch; %v - %v", testCase.weather, weth)
				}
				for i := range weth.List {
					readDay := weth.List[i]
					putDay := testCase.weather.List[i]

					if !(readDay.Date.Equal(putDay.Date)) || (readDay.Desc != putDay.Desc) || (readDay.Feels != putDay.Feels) || (readDay.Temp != putDay.Temp) {
						t.Errorf("Mismatch: %v - %v", readDay, putDay)
					}
				}
			}
		case LatLon:
			{
				if (weth.Lat != testCase.weather.Lat) || (weth.Lon != testCase.weather.Lon) || (weth.Type != testCase.weather.Type) || (len(weth.List) != len(testCase.weather.List)) {
					t.Errorf("Mismatch; %v - %v", testCase.weather, weth)
				}
				for i := range weth.List {
					readDay := weth.List[i]
					putDay := testCase.weather.List[i]

					if !(readDay.Date.Equal(putDay.Date)) || (readDay.Desc != putDay.Desc) || (readDay.Feels != putDay.Feels) || (readDay.Temp != putDay.Temp) {
						t.Errorf("Mismatch: %v - %v", readDay, putDay)
					}
				}
			}
		default:
			{
				t.Errorf("unexpected type %T", tp)
			}
		}

	}

}

func TestValidateCacheEntry(t *testing.T) {

	now := time.Now()

	testCases := []struct {
		queryType QueryType
		date      time.Time
		expect    bool
	}{
		{QueryType(Current), now, true},
		{QueryType(Current), now.Add(-time.Minute * 10).Add(time.Second), true},
		{QueryType(Current), now.Add(-time.Minute * 10).Add(-time.Second), false},
		{QueryType(Forecast), now, true},
		{QueryType(Forecast), now.Add(-time.Hour * 3).Add(time.Second), true},
		{QueryType(Forecast), now.Add(-time.Hour * 3).Add(-time.Second), false},
	}

	for _, testCase := range testCases {
		res := ValidateCacheEntry(testCase.queryType, testCase.date)
		if testCase.expect != res {
			t.Errorf("Mismatch. %v - got %v; want %v", testCase.date, res, testCase.expect)
		}
	}

}
