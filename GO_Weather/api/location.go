package api

import (
	"fmt"
	"net/url"
	"strings"
)

type Location interface {
	CacheKey() (string, error)
	QueryParam() (string, error)
}

type CityCountry struct {
	City    string
	Country string
}

func (cc CityCountry) QueryParam() (string, error) {

	if cc.City != "" {
		if cc.Country != "" {
			return fmt.Sprintf("q=%s,%s", url.QueryEscape(cc.City), cc.Country), nil
		}
		return fmt.Sprintf("q=%s", url.QueryEscape(cc.City)), nil
	}

	if cc.Country != "" {
		return fmt.Sprintf("q=%s", cc.Country), nil
	}

	return "", fmt.Errorf("missing values")

}

func (cc CityCountry) CacheKey() (string, error) {

	cacheKey := ""
	city := strings.ToLower(cc.City)
	country := strings.ToLower(cc.Country)
	countryCode := CountryCodes[country]

	if city != "" {
		cacheKey += city
		if country != "" {
			cacheKey += "," + countryCode
		}
	} else if country != "" {
		cacheKey += countryCode
	} else {
		return "", fmt.Errorf("invalid or missing data")
	}

	return cacheKey, nil
}

type LatLon struct {
	Lat float64
	Lon float64
}

func (ll LatLon) QueryParam() (string, error) {

	if ll.Lat != 0 || ll.Lon != 0 {
		return fmt.Sprintf("lat=%f&lon=%f", ll.Lat, ll.Lon), nil
	}

	return "", fmt.Errorf("missing values")

}

func (ll LatLon) CacheKey() (string, error) {

	if ll.Lat != 0 || ll.Lon != 0 {
		return fmt.Sprintf("%.2f-%.2f", ll.Lat, ll.Lon), nil
	}

	return "", fmt.Errorf("invalid or missing data")

}
