package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func main() {

	const url string = "http://api.openweathermap.org/data/2.5/weather"
	var city string

	file_data, file_err := os.ReadFile(".apikey")
	if file_err != nil {
		fmt.Println("Error reading the api key", file_err)
	}

	api_key := strings.TrimSpace(string(file_data))

	flag.StringVar(&city, "city", "Tel Aviv", "City name")
	flag.StringVar(&city, "c", "Tel Aviv", "Shorthand for -city")
	flag.Parse()

	weather_query := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", url, strings.ReplaceAll(city, " ", "+"), api_key)

	resp, err := http.Get(weather_query)
	if err != nil {
		fmt.Println("Error with the request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
        	fmt.Printf("API error: %s (status code %d)\n", resp.Status, resp.StatusCode)
        	os.Exit(1)
	}

	var data WeatherResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	fmt.Printf("%s: %.1f°C, %s\n", data.Name, data.Main.Temp, strings.Title(data.Weather[0].Description))
}
