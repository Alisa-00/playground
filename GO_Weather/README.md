# Weather API Client in Go

This project is a simple CLI Go application designed to fetch and display current weather data for a given city using the [OpenWeatherMap API](https://openweathermap.org/api).

---

## Overview

This is a small Go learning project, it demonstrates how to:

- Organize Go code into packages and modules
- Make HTTP requests and handle JSON responses
- Basic file handling and environment configuration
- Manage configuration files (such as API keys) safely
- Handle errors in Go functions

---

## How to Use

1. Clone the repository.
2. Create a `.apikey` file in the `api` folder containing your OpenWeatherMap API key.
3. Run the program with:

   ```bash
   go run main.go -c "City name"
