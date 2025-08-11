package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const baseUrl string = "https://api.openweathermap.org/data/2.5"
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
	Lat     float64
	Lon     float64
	List    []WeatherDay
}

type WeatherDay struct {
	Temp  float64
	Feels float64
	Desc  string
	Date  time.Time
}

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp  float64 `json:"temp"`
		Feels float64 `json:"feels_like"`
	} `json:"main"`
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
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
		Coord   struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
	} `json:"city"`
}

func getQueryString(loc Location, queryType QueryType, units string) (string, error) {

	queryUrl := baseUrl + "/"

	switch queryType {
	case QueryType(Current):
		{
			queryUrl += weatherEndpoint
		}
	case QueryType(Forecast):
		{
			queryUrl += forecastEndpoint
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

	queryParams, err := loc.QueryParam()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s?appid=%s&units=%s&%s", queryUrl, apiKey, units, queryParams), nil
}

func CreateLocation(city string, country string, lat float64, lon float64) (Location, error) {

	if lat != 0 || lon != 0 {
		return LatLon{Lat: lat, Lon: lon}, nil
	}

	if city != "" {
		if country != "" {
			return CityCountry{City: city, Country: country}, nil
		}
		return CityCountry{City: city}, nil
	}

	if country != "" {
		return CityCountry{Country: country}, nil
	}

	return CityCountry{}, fmt.Errorf("invalid input. city, country or latitude and longitude have to be valid")
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
		return data, fmt.Errorf("error executing query: %w. %s", err, queryUrl)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return data, fmt.Errorf("bad status code: %s (status code %d)", resp.Status, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("error decoding JSON response: %w. %s", err, queryUrl)
	}

	return data, nil
}

func GetCurrentWeather(loc Location, units string) (Weather, error) {

	weatherData := Weather{Type: QueryType(Current)}
	query, err := getQueryString(loc, weatherData.Type, units)
	if err != nil {
		return weatherData, err
	}

	data, err := fetchWeather[WeatherResponse](query)
	if err != nil {
		return weatherData, err
	}

	if len(data.Weather) == 0 {
		return weatherData, fmt.Errorf("no weather data found")
	}

	weatherData.City = data.Name
	weatherData.Country = data.Sys.Country
	weatherData.Lat = data.Coord.Lat
	weatherData.Lon = data.Coord.Lon

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
	query, err := getQueryString(loc, weatherData.Type, units)
	if err != nil {
		return weatherData, err
	}

	data, err := fetchWeather[ForecastResponse](query)
	if err != nil {
		return weatherData, err
	}

	if len(data.List) == 0 {
		return weatherData, fmt.Errorf("no forecast data found")
	}

	weatherData.City = data.City.Name
	weatherData.Country = data.City.Country
	weatherData.Lat = data.City.Coord.Lat
	weatherData.Lon = data.City.Coord.Lon

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
