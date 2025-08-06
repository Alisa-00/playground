package main

import (
	"flag"
	"fmt"
	"weather/api"
	"weather/parse"
)

func main() {

	// read flags
	var city, country, units string
	var lat, lon float64
	var forecast bool
	flag.StringVar(&city, "city", "", "City name")
	flag.StringVar(&city, "c", "", "Shorthand for -city")
	flag.StringVar(&country, "country", "", "Country code")
	flag.StringVar(&country, "C", "", "Shorthand for country code")
	flag.StringVar(&units, "units", "metric", "Units for displaying")
	flag.StringVar(&units, "u", "metric", "Shorthand for units")
	flag.BoolVar(&forecast, "f", false, "Forecast")
	flag.Float64Var(&lat, "lat", 0, "Latitude")
	flag.Float64Var(&lon, "lon", 0, "Longitude")
	flag.Parse()

	// parse and validate inputs
	units, err := parse.ParseUnits(units)
	if err != nil {
		fmt.Println(err)
		return
	}

	char, err := parse.GetChar(units)
	if err != nil {
		fmt.Println(err)
		return
	}

	location, err := api.CreateLocation(city, country, lat, lon)
	if err != nil {
		fmt.Println(err)
		return
	}

	// query API and display results
	var weather api.Weather
	if forecast {
		weath, err := api.GetForecast(location, units)
		if err != nil {
			fmt.Println(err)
			return
		}
		weather = weath
	} else {
		weath, err := api.GetWeather(location, units)
		if err != nil {
			fmt.Println(err)
			return
		}
		weather = weath
	}

	display(weather, char)
}

func display(weather api.Weather, char rune) {
	fmt.Printf("%s,%s\n", weather.City, weather.Country)
	for _, day := range weather.List {
		fmt.Printf("%s: %.1f°%c, %s. It feels like %.1f°%c\n", day.Date, day.Temp, char, day.Desc, day.Feels, char)
	}
}
