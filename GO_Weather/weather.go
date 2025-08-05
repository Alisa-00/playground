package main

import (
	"flag"
	"fmt"
	"weather/api"
)

func main() {

	// read flags
	var city, country string
	var lat, lon float64
	flag.StringVar(&city, "city", "", "City name")
	flag.StringVar(&city, "c", "", "Shorthand for -city")
	flag.StringVar(&country, "country", "", "Country code")
	flag.StringVar(&country, "C", "", "Shorthand for country code")
	flag.Float64Var(&lat, "lat", 0, "Latitude")
	flag.Float64Var(&lon, "lon", 0, "Longitude")
	flag.Parse()

	location, err := getLocation(city, country, lat, lon)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get weather
	weather, err := api.GetWeather(location)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print on screen
	display(weather)
}

func getLocation(city string, country string, lat float64, lon float64) (api.Location, error) {

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

func display(weather api.Weather) {
	fmt.Printf("%s: %.1f°C, %s\n", weather.City, weather.Temp, weather.Desc)
}
