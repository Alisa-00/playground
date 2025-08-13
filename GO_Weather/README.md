# Weather CLI

A fast, modular, and cache-enabled command-line weather tool written in Go.  
Uses [OpenWeatherMap API](https://openweathermap.org/api) to fetch current conditions and forecasts for any location using either **city**, **country**, **city/country** or **latitude/longitude** input, with caching of recent queries to minimize API calls and improve performance.

---

## Features

- **Query Modes**
  - Search by city and country
  - Search by geographic coordinates
  - Support for both **current weather** and **forecast** queries.
  - Support for displaying different temperature units

- **Caching**
  - Reduces redundant API calls by storing responses in a JSON cache file.
  - Automatic cache invalidation after a time limit.

## Usage

```bash
# Current weather by city and country
go run . -c Paris -C France
```
Paris,FR
Aug 13 12:22:00: 28.7°C, clear sky. It feels like 29.5°C

```bash
# Current weather by coordinates
go run . -lat 48.8566 -lon 2.3522
```
Paris,FR
Aug 13 12:23:46: 28.6°C, clear sky. It feels like 29.4°C

```bash
# Forecast by city and country
go run . -c "New York" -C America -f -u F
```
New York,US
Aug 13 15:00:00: 24.0°F, light rain. It feels like 24.7°F
Aug 13 18:00:00: 28.4°F, scattered clouds. It feels like 31.0°F
Aug 13 21:00:00: 32.1°F, scattered clouds. It feels like 35.9°F
Aug 14 00:00:00: 28.2°F, light rain. It feels like 30.6°F
Aug 14 03:00:00: 25.2°F, moderate rain. It feels like 25.7°F
Aug 14 06:00:00: 25.4°F, light rain. It feels like 26.1°F
Aug 14 09:00:00: 25.1°F, overcast clouds. It feels like 25.8°F
Aug 14 12:00:00: 25.0°F, overcast clouds. It feels like 25.7°F
Aug 14 15:00:00: 25.6°F, overcast clouds. It feels like 26.1°F
Aug 14 18:00:00: 25.6°F, overcast clouds. It feels like 26.1°F
Aug 14 21:00:00: 26.2°F, light rain. It feels like 26.2°F
Aug 15 00:00:00: 29.2°F, overcast clouds. It feels like 30.6°F
Aug 15 03:00:00: 27.0°F, light rain. It feels like 28.7°F
Aug 15 06:00:00: 25.2°F, light rain. It feels like 25.8°F
Aug 15 09:00:00: 25.0°F, light rain. It feels like 25.6°F
Aug 15 12:00:00: 24.4°F, clear sky. It feels like 25.0°F
Aug 15 15:00:00: 24.6°F, clear sky. It feels like 24.9°F
Aug 15 18:00:00: 28.3°F, clear sky. It feels like 29.5°F
Aug 15 21:00:00: 30.2°F, clear sky. It feels like 31.4°F
Aug 16 00:00:00: 28.6°F, clear sky. It feels like 29.9°F
Aug 16 03:00:00: 26.2°F, few clouds. It feels like 26.2°F
Aug 16 06:00:00: 25.2°F, broken clouds. It feels like 25.5°F
Aug 16 09:00:00: 23.8°F, scattered clouds. It feels like 24.2°F
Aug 16 12:00:00: 23.3°F, scattered clouds. It feels like 23.6°F
Aug 16 15:00:00: 23.9°F, scattered clouds. It feels like 24.1°F
Aug 16 18:00:00: 24.8°F, overcast clouds. It feels like 25.0°F
Aug 16 21:00:00: 27.3°F, broken clouds. It feels like 27.9°F
Aug 17 00:00:00: 27.2°F, clear sky. It feels like 27.8°F
Aug 17 03:00:00: 24.8°F, scattered clouds. It feels like 25.0°F
Aug 17 06:00:00: 24.3°F, few clouds. It feels like 24.6°F
Aug 17 09:00:00: 23.7°F, few clouds. It feels like 24.0°F
Aug 17 12:00:00: 23.0°F, clear sky. It feels like 23.3°F
Aug 17 15:00:00: 24.0°F, clear sky. It feels like 24.2°F
Aug 17 18:00:00: 29.0°F, clear sky. It feels like 30.0°F
Aug 17 21:00:00: 32.7°F, clear sky. It feels like 34.2°F
Aug 18 00:00:00: 31.2°F, light rain. It feels like 32.7°F
Aug 18 03:00:00: 28.1°F, light rain. It feels like 30.1°F
Aug 18 06:00:00: 27.2°F, scattered clouds. It feels like 29.2°F
Aug 18 09:00:00: 26.0°F, light rain. It feels like 26.0°F
Aug 18 12:00:00: 24.8°F, light rain. It feels like 25.3°F