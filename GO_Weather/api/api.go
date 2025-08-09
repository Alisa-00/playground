package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseUrl string = "http://api.openweathermap.org/data/2.5"
const weatherEndpoint string = "weather"
const forecastEndpoint string = "forecast"
const DateFormat = time.Stamp
const apiKeyEnvName string = "OPENWEATHER_API_KEY"

type QueryType int

const (
	Current QueryType = iota
	Forecast
)

type Weather struct {
	Type    QueryType
	City    string
	Country string
	List    []WeatherDay
}

type WeatherDay struct {
	Temp  float64
	Feels float64
	Desc  string
	Date  time.Time
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

func (loc Location) getQueryString(queryType QueryType, units string) (string, error) {

	url := baseUrl + "/"

	switch queryType {
	case QueryType(Current):
		{
			url += weatherEndpoint
		}
	case QueryType(Forecast):
		{
			url += forecastEndpoint
		}
	default:
		{
			return "", fmt.Errorf("invalid query type")
		}
	}

	apiKey, err := getApiKey()
	if err != nil {
		return "", err
	}
	baseQuery := fmt.Sprintf("%s?appid=%s&units=%s", url, apiKey, units)

	if loc.Latitude != 0 || loc.Longitude != 0 {
		return fmt.Sprintf("%s&lat=%f&lon=%f", loc.Latitude, loc.Longitude), nil
	}

	if loc.City != "" {
		if loc.Country != "" {
			return fmt.Sprintf("%s&q=%s,%s", baseQuery, strings.ReplaceAll(loc.City, " ", "+"), loc.Country), nil
		}
		return fmt.Sprintf("%s&q=%s", baseQuery, strings.ReplaceAll(loc.City, " ", "+")), nil
	}

	if loc.Country != "" {
		return fmt.Sprintf("%s&q=%s", baseQuery, loc.Country), nil
	}

	return "", fmt.Errorf("missing values, query was not constructed")
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

func fetchWeather[T any](queryUrl string) (T, error) {
	var data T
	resp, err := http.Get(queryUrl)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return data, fmt.Errorf("bad status code: %s (status code %d)", resp.Status, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func GetCurrentWeather(loc Location, units string) (Weather, error) {

	weatherData := Weather{Type: QueryType(Current)}
	query, err := loc.getQueryString(weatherData.Type, units)
	if err != nil {
		return weatherData, nil
	}

	data, err := fetchWeather[WeatherResponse](query)
	if err != nil {
		return weatherData, nil
	}

	if len(data.Weather) == 0 {
		return weatherData, fmt.Errorf("no weather data found")
	}

	weatherData.City = data.Name
	weatherData.Country = data.Sys.Country
	weatherData.List = []WeatherDay{
		{
			Temp:  data.Main.Temp,
			Feels: data.Main.Feels,
			Desc:  data.Weather[0].Description,
			Date:  time.Unix(data.Date, 0),
		},
	}

	return weatherData, nil
}

func GetForecast(loc Location, units string) (Weather, error) {

	weatherData := Weather{Type: QueryType(Forecast)}
	query, err := loc.getQueryString(weatherData.Type, units)
	if err != nil {
		return weatherData, err
	}

	data, err := fetchWeather[ForecastResponse](query)
	if err != nil {
		return weatherData, nil
	}

	if len(data.List) == 0 {
		return weatherData, fmt.Errorf("no forecast data found")
	}

	weatherData.City = data.City.Name
	weatherData.Country = data.City.Country

	weatherList := make([]WeatherDay, 0)
	for _, data := range data.List {
		day := WeatherDay{
			Temp:  data.Main.Temp,
			Feels: data.Main.Feels,
			Desc:  data.Weather[0].Description,
			Date:  time.Unix(data.Date, 0),
		}
		weatherList = append(weatherList, day)
	}
	weatherData.List = weatherList

	return weatherData, nil
}
