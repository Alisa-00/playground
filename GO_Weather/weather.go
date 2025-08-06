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
	flag.StringVar(&city, "city", "", "City name")
	flag.StringVar(&city, "c", "", "Shorthand for -city")
	flag.StringVar(&country, "country", "", "Country code")
	flag.StringVar(&country, "C", "", "Shorthand for country code")
	flag.StringVar(&units, "units", "metric", "Units for displaying")
	flag.StringVar(&units, "u", "metric", "Shorthand for units")
	flag.Float64Var(&lat, "lat", 0, "Latitude")
	flag.Float64Var(&lon, "lon", 0, "Longitude")
	flag.Parse()

	units, err := parse.ParseUnits(units)
	if err != nil {
		fmt.Println(err)
		return
	}
	location, err := api.CreateLocation(city, country, lat, lon)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get weather
	weather, err := api.GetWeather(location, units)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print on screen
	char, err := parse.GetChar(units)
	if err != nil {
		fmt.Println(err)
		return
	}
	display(weather, char)
}

func display(weather api.Weather, char rune) {
	fmt.Printf("%s,%s: %.1f°%c, %s. It feels like %.1f°%c\n", weather.City, weather.Country, weather.Temp, char, weather.Desc, weather.Feels, char)
}
