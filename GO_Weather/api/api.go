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

func GetWeather(city string) (Weather, error) {

	const url string = URL
	apiKey, err := getApiKey()
	if err != nil {
		return Weather{}, fmt.Errorf("API Key not found")
	}
	query := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", url, strings.ReplaceAll(city, " ", "+"), apiKey)

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
