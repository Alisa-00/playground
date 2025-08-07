package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseUrl string = "http://api.openweathermap.org/data/2.5"
const weatherEndpoint string = "weather"
const forecastEndpoint string = "forecast"
const dateFormat = time.Stamp
const apiKeyEnvName string = "OPENWEATHER_API_KEY"

type Weather struct {
	City    string
	Country string
	List    []weatherDay
}

type weatherDay struct {
	Temp  float64
	Feels float64
	Desc  string
	Date  string
}

type Location struct {
	City      string
	Country   string
	Latitude  float64
	Longitude float64
}

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp  float64 `json:"temp"`
		Feels float64 `json:"feels_like"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Date int64 `json:"dt"`
}

type ForecastResponse struct {
	List []struct {
		Main struct {
			Temp  float64 `json:"temp"`
			Feels float64 `json:"feels_like"`
		} `json:"main"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
		Date int64 `json:"dt"`
	} `json:"list"`
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
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

func CreateLocation(city string, country string, lat float64, lon float64) (Location, error) {

	if lat != 0 || lon != 0 {
		return Location{Latitude: lat, Longitude: lon}, nil
	}

	if city != "" {
		if country != "" {
			return Location{City: city, Country: country}, nil
		}
		return Location{City: city}, nil
	}

	if country != "" {
		return Location{Country: country}, nil
	}

	return Location{}, fmt.Errorf("invalid input. city, country or latitude and longitude have to be valid")
}

func getApiKey() (string, error) {
	apiKey := os.Getenv(apiKeyEnvName)

	if apiKey == "" {
		return "", fmt.Errorf("api key unavailable. $%s is not set", apiKeyEnvName)
	}

	return apiKey, nil
}

func GetWeather(loc Location, units string) (Weather, error) {

	apiKey, err := getApiKey()
	if err != nil {
		return Weather{}, err
	}

	locationParams := loc.getQuerySubstring()
	query := fmt.Sprintf("%s/%s?%s&appid=%s&units=%s", baseUrl, weatherEndpoint, locationParams, apiKey, units)

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
		City:    data.Name,
		Country: data.Sys.Country,
		List: []weatherDay{
			{
				Temp:  data.Main.Temp,
				Feels: data.Main.Feels,
				Desc:  data.Weather[0].Description,
				Date:  time.Unix(data.Date, 0).Format(dateFormat),
			},
		},
	}, nil
}

func GetForecast(loc Location, units string) (Weather, error) {

	apiKey, err := getApiKey()
	if err != nil {
		return Weather{}, err
	}

	locationParams := loc.getQuerySubstring()
	query := fmt.Sprintf("%s/%s?%s&appid=%s&units=%s", baseUrl, forecastEndpoint, locationParams, apiKey, units)

	resp, err := http.Get(query)
	if err != nil {
		return Weather{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Weather{}, fmt.Errorf("bad status code: %s (status code %d)", resp.Status, resp.StatusCode)
	}

	var data ForecastResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return Weather{}, fmt.Errorf("invalid JSON: %w", err)
	}

	if len(data.List) == 0 {
		return Weather{}, fmt.Errorf("no forecast data found")
	}

	weatherList := make([]weatherDay, 0)
	for _, data := range data.List {
		day := weatherDay{
			Temp:  data.Main.Temp,
			Feels: data.Main.Feels,
			Desc:  data.Weather[0].Description,
			Date:  time.Unix(data.Date, 0).Format(dateFormat),
		}
		weatherList = append(weatherList, day)
	}

	return Weather{
		City:    data.City.Name,
		Country: data.City.Country,
		List:    weatherList,
	}, nil
}
