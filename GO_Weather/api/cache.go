package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	return &cache, nil
}

func (cache Cache) SaveCacheFile() error {

	// invalidate old entries
	cache.invalidateCache()

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

func normalizeCoord(val float64) string {
	return fmt.Sprintf("%.2f", val)
}

func getCacheKeyType(key string, queryType QueryType) (string, error) {

	switch queryType {
	case QueryType(Current):
		{
			key += "_" + "current"
		}
	case QueryType(Forecast):
		{
			key += "_" + "forecast"
		}
	default:
		{
			return "", fmt.Errorf("invalid query type")
		}
	}

	return key, nil
}

func getCacheKeyCC(city string, country string, queryType QueryType) (string, error) {

	cacheKey := ""

	if city != "" {
		cacheKey += city
		if country != "" {
			cacheKey += "," + country
		}
	} else if country != "" {
		cacheKey += country
	} else {
		return "", fmt.Errorf("invalid or missing data")
	}

	cacheKey, err := getCacheKeyType(cacheKey, queryType)
	if err != nil {
		return "", err
	}

	return cacheKey, nil

}

func getCacheKeyLL(lat float64, lon float64, queryType QueryType) (string, error) {

	cacheKey := ""

	if lat != 0 || lon != 0 {
		cacheKey += normalizeCoord(lat) + "-" + normalizeCoord(lon)
	} else {
		return "", fmt.Errorf("invalid or missing data")
	}

	cacheKey, err := getCacheKeyType(cacheKey, queryType)
	if err != nil {
		return "", err
	}

	return cacheKey, nil

}

func (cache Cache) invalidateCache() {
	for key, weather := range cache {
		valid := ValidateCacheEntry(weather.Type, weather.List[0].Date)
		if !valid {
			delete(cache, key)
		}
	}
}

func ValidateCacheEntry(queryType QueryType, date time.Time) bool {

	hours := time.Since(date).Abs().Hours()

	switch queryType {
	case QueryType(Current):
		{
			if hours <= hoursInvalidateCurrent {
				return true
			}
		}
	case QueryType(Forecast):
		{
			if hours <= hoursInvalidateForecast {
				return true
			}
		}
	}

	return false
}

func (cache Cache) Put(weather Weather) error {

	ccCacheKey, err := getCacheKeyCC(weather.City, weather.Country, weather.Type)
	if err != nil {
		return err
	}
	llCacheKey, err := getCacheKeyLL(weather.Lat, weather.Lon, weather.Type)
	if err != nil {
		return err
	}

	cache[ccCacheKey] = weather
	cache[llCacheKey] = weather

	return nil
}

func (cache Cache) ReadCC(loc Location, queryType string) (Weather, error) {

	weather := Weather{}
	switch queryType {
	case "current":
		{
			cacheKey, err := getCacheKeyCC(loc.City, loc.Country, QueryType(Current))
			if err != nil {
				return weather, err
			}
			weather = cache[cacheKey]

		}
	case "forecast":
		{
			cacheKey, err := getCacheKeyCC(loc.City, loc.Country, QueryType(Forecast))
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

func (cache Cache) ReadLL(loc Location, queryType string) (Weather, error) {

	weather := Weather{}
	switch queryType {
	case "current":
		{
			cacheKey, err := getCacheKeyLL(loc.Latitude, loc.Longitude, QueryType(Current))
			if err != nil {
				return weather, err
			}
			weather = cache[cacheKey]

		}
	case "forecast":
		{
			cacheKey, err := getCacheKeyLL(loc.Latitude, loc.Longitude, QueryType(Forecast))
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
