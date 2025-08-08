package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const cacheFile string = ".cache"
const tempCacheFile string = ".cache_temp"
const configDir string = ".config"
const appDir string = "go_weather_cli"
const hoursInvalidateCurrent float64 = 0.167
const hoursInvalidateForecast float64 = 3

type Cache map[string]Weather

func LoadCacheFile() (*Cache, error) {

	cache := make(Cache)

	// get home dir and paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &cache, err
	}

	configPath := filepath.Join(homeDir, configDir, appDir)
	cacheFile := filepath.Join(configPath, cacheFile)

	// read from cache file
	bytes, err := os.ReadFile(cacheFile)
	if err != nil {
		cache.SaveCacheFile()
		return &cache, err
	}

	// load cache data into map
	err = json.Unmarshal(bytes, &cache)
	if err != nil {
		return &cache, err
	}

	// invalidate outdated entries

	for key, weather := range cache {
		if strings.Contains(key, "current") {
			weatherDate := weather.List[0].Date
			hours := time.Since(weatherDate).Abs().Hours()
			if hours > hoursInvalidateCurrent {
				delete(cache, key)
			}
		}
		if strings.Contains(key, "forecast") {
			weatherDate := weather.List[0].Date
			hours := time.Since(weatherDate).Abs().Hours()
			if hours > hoursInvalidateForecast {
				delete(cache, key)
			}
		}
	}

	return &cache, nil
}

func (cache Cache) SaveCacheFile() error {

	// get home dir and paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, configDir, appDir)
	cacheFile := filepath.Join(configPath, cacheFile)
	tempCacheFile := filepath.Join(configPath, tempCacheFile)

	// create dir if needed
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	// write into temp file
	err = os.WriteFile(tempCacheFile, jsonBytes, 0644)
	if err != nil {
		return err
	}

	// rename to permanent cache file
	err = os.Rename(tempCacheFile, cacheFile)
	if err != nil {
		return err
	}

	return nil
}

func getCacheKey(city string, queryType QueryType) (string, error) {

	cacheKey := strings.ToLower(city) + "_"

	switch queryType {
	case QueryType(Current):
		{
			cacheKey = cacheKey + "current"
		}
	case QueryType(Forecast):
		{
			cacheKey = cacheKey + "forecast"
		}
	default:
		{
			return "", fmt.Errorf("invalid query type")
		}
	}

	return cacheKey, nil
}

func (cache Cache) Put(weather Weather) error {

	entry, err := getCacheKey(weather.City, weather.Type)
	if err != nil {
		return err
	}
	cache[entry] = weather

	return nil
}

func (cache Cache) ReadCC(city string, queryType string) (Weather, error) {

	weather := Weather{}
	switch queryType {
	case "current":
		{
			cacheKey, err := getCacheKey(city, QueryType(Current))
			if err != nil {
				return weather, err
			}
			weather = cache[cacheKey]

		}
	case "forecast":
		{
			cacheKey, err := getCacheKey(city, QueryType(Forecast))
			if err != nil {
				return weather, err
			}
			weather = cache[cacheKey]
		}
	default:
		{
			return weather, fmt.Errorf("invalid query type")
		}
	}

	return weather, nil
}

func (cache Cache) ReadLL(lat float64, lon float64, queryType string) (Weather, error) {
	weather := Weather{}
	return weather, nil
}
