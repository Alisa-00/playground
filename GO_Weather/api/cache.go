package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheFileName string = ".cache"
const tempCacheFileName string = ".cache_temp"
const configDir string = ".config"
const appDir string = "go_weather_cli"
const invalidateCurrent = 10 * time.Minute
const invalidateForecast = 3 * time.Hour

type Cache map[string]Weather

func LoadCacheFile() (*Cache, error) {

	cache := make(Cache)

	// get home dir and paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &cache, err
	}

	configPath := filepath.Join(homeDir, configDir, appDir)
	cacheFile := filepath.Join(configPath, cacheFileName)

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
	cacheFile := filepath.Join(configPath, cacheFileName)
	tempCacheFile := filepath.Join(configPath, tempCacheFileName)

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

func getCacheKey(loc Location, queryType QueryType) (string, error) {

	cacheKey, err := loc.CacheKey()
	if err != nil {
		return "", err
	}

	cacheKey, err = getCacheKeyType(cacheKey, queryType)
	if err != nil {
		return "", err
	}

	return cacheKey, nil

}

func (cache Cache) invalidateCache() {
	for key, weather := range cache {
		if len(weather.List) == 0 {
			delete(cache, key)
			continue
		}
		valid := ValidateCacheEntry(weather.Type, weather.List[0].Date)
		if !valid {
			delete(cache, key)
		}
	}
}

func ValidateCacheEntry(queryType QueryType, date time.Time) bool {

	timeSpan := time.Since(date)

	switch queryType {
	case QueryType(Current):
		{
			if timeSpan <= invalidateCurrent {
				return true
			}
		}
	case QueryType(Forecast):
		{
			if timeSpan <= invalidateForecast {
				return true
			}
		}
	}

	return false
}

func (cache Cache) Put(loc Location, weather Weather) error {

	ccCacheKey, err := getCacheKey(loc, weather.Type)
	if err == nil {
		cache[ccCacheKey] = weather
	}

	return err
}

func (cache Cache) Read(loc Location, queryType string) (Weather, error) {

	weather := Weather{}
	switch queryType {
	case "current":
		{
			cacheKey, err := getCacheKey(loc, QueryType(Current))
			if err != nil {
				return weather, err
			}
			weather = cache[cacheKey]

		}
	case "forecast":
		{
			cacheKey, err := getCacheKey(loc, QueryType(Forecast))
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
