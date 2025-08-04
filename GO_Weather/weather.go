package main

import (
	"flag"
	"fmt"
	"weather/api"
)

func main() {

	// read flags
	var city string
	flag.StringVar(&city, "city", "Tel Aviv", "City name")
	flag.StringVar(&city, "c", "Tel Aviv", "Shorthand for -city")
	flag.Parse()

	// get weather
	weather, err := api.GetWeather(city)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print on screen
	display(weather)
}

func display(weather api.Weather) {

	fmt.Printf("%s: %.1f°C, %s\n", weather.City, weather.Temp, weather.Desc)
}
