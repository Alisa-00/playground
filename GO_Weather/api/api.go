package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const URL string = "http://api.openweathermap.org/data/2.5/weather"

type Weather struct {
	City string
	Temp float64
	Desc string
}

type Location struct {
	City      string
	Country   string
	Latitude  float64
	Longitude float64
}

func (loc Location) getQuerySubstring() string {
	if loc.Latitude != 0 || loc.Longitude != 0 {
		return fmt.Sprintf("lat=%f&lon=%f", loc.Latitude, loc.Longitude)
	}

	if loc.City != "" {
		if loc.Country != "" {
			return fmt.Sprintf("q=%s,%s", strings.ReplaceAll(loc.City, " ", "+"), loc.Country)
		}
		return fmt.Sprintf("q=%s", strings.ReplaceAll(loc.City, " ", "+"))
	}

	if loc.Country != "" {
		return fmt.Sprintf("q=%s", loc.Country)
	}

	return ""
}

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func getConfigPath(filename string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, filename)
}

func getApiKey() (string, error) {
	apiKeyPath := getConfigPath(".apikey")
	data, err := os.ReadFile(apiKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading the api key: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func GetWeather(loc Location) (Weather, error) {

	const url string = URL
	apiKey, err := getApiKey()
	if err != nil {
		return Weather{}, fmt.Errorf("API Key not found")
	}
	locationParams := loc.getQuerySubstring()
	query := fmt.Sprintf("%s?%s&appid=%s&units=metric", url, locationParams, apiKey)

	resp, err := http.Get(query)
	if err != nil {
		return Weather{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Weather{}, fmt.Errorf("bad status code: %s (status code %d)", resp.Status, resp.StatusCode)
	}

	var data WeatherResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return Weather{}, fmt.Errorf("invalid JSON: %w", err)
	}

	if len(data.Weather) == 0 {
		return Weather{}, fmt.Errorf("no weather data found")
	}

	return Weather{
		data.Name,
		data.Main.Temp,
		data.Weather[0].Description,
	}, nil
}
