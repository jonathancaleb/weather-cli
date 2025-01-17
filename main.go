package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=d28a46779f2d40d581375333251701&q=Kampala&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Failed to fetch weather data")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current := weather.Location, weather.Current

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	// Safely access hourly forecast data
	if len(weather.Forecast.Forecastday) > 0 {
		hours := weather.Forecast.Forecastday[0].Hour
		for _, hour := range hours {
			date := time.Unix(hour.TimeEpoch, 0)
			fmt.Printf(
				"%s - %.0fC, %.0f%% chance of rain, %s\n",
				date.Format("15:04"),
				hour.TempC,
				hour.ChanceOfRain,
				hour.Condition.Text,
			)
		}
	} else {
		fmt.Println("No forecast data available.")
	}
}
