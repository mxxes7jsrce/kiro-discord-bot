package bot

import (
	"fmt"
	"math/rand"
)

// WeatherCondition represents a simple weather report.
type WeatherCondition struct {
	Condition   string
	TemperatureC int
	Humidity    int
	Description string
}

var weatherConditions = []WeatherCondition{
	{"Sunny", 28, 40, "Clear skies and bright sunshine."},
	{"Cloudy", 18, 60, "Overcast with thick clouds."},
	{"Rainy", 14, 85, "Light to moderate rainfall expected."},
	{"Stormy", 10, 90, "Thunderstorms with heavy rain."},
	{"Snowy", -2, 75, "Light snowfall throughout the day."},
	{"Windy", 20, 50, "Strong gusts of wind, stay safe!"},
	{"Foggy", 12, 95, "Dense fog reducing visibility."},
	{"Partly Cloudy", 22, 55, "Mix of sun and clouds."},
}

// GetWeather returns a simulated weather report for a given location.
func GetWeather(location string) (*WeatherCondition, error) {
	if location == "" {
		return nil, fmt.Errorf("location cannot be empty")
	}
	// Simulate weather by seeding from location string
	seed := int64(0)
	for _, c := range location {
		seed += int64(c)
	}
	r := rand.New(rand.NewSource(seed))
	condition := weatherConditions[r.Intn(len(weatherConditions))]
	return &condition, nil
}

// FormatWeatherReport formats a WeatherCondition into a Discord message.
func FormatWeatherReport(location string, w *WeatherCondition) string {
	return fmt.Sprintf(
		"🌍 **Weather for %s**\n🌤 Condition: %s\n🌡 Temperature: %d°C\n💧 Humidity: %d%%\n📝 %s",
		location, w.Condition, w.TemperatureC, w.Humidity, w.Description,
	)
}
