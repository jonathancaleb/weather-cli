package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
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
	now := time.Now()

	// Display current weather
	conditionEmoji := getConditionEmoji(current.Condition.Text)
	fmt.Printf(
		"%s, %s: %.0f°C %s %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
		conditionEmoji,
	)

	// Check for rain later in the day
	if len(weather.Forecast.Forecastday) > 0 {
		hours := weather.Forecast.Forecastday[0].Hour
		rainExpected := false
		daySummary := "Sunny Day 🌞"

		for _, hour := range hours {
			hourTime := time.Unix(hour.TimeEpoch, 0)

			// Skip past hours
			if hourTime.Before(now) {
				continue
			}

			// Check if rain is expected later
			if hour.ChanceOfRain > 40 {
				rainExpected = true
				daySummary = "Rainy Day 🌧"
			}
		}

		if rainExpected {
			color.Blue("⛈ Rain expected later. Forecast for the day: %s", daySummary)
		} else {
			color.Green("🌞 No rain expected. Forecast for the day: %s", daySummary)
		}
	} else {
		fmt.Println("No forecast data available.")
	}
}

// getConditionEmoji maps weather conditions to emojis
func getConditionEmoji(condition string) string {
	switch condition {
	case "Sunny", "Clear":
		return "☀️"
	case "Rain", "Light rain", "Showers":
		return "🌧"
	case "Cloudy", "Overcast":
		return "☁️"
	case "Snow":
		return "❄️"
	case "Thunderstorm":
		return "⛈"
	default:
		return "🌡"
	}
}
