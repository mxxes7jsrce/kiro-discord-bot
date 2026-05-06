package bot

import (
	"strings"
	"testing"
)

func TestGetWeather_EmptyLocation(t *testing.T) {
	_, err := GetWeather("")
	if err == nil {
		t.Fatal("expected error for empty location, got nil")
	}
}

func TestGetWeather_ValidLocation(t *testing.T) {
	w, err := GetWeather("London")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected weather condition, got nil")
	}
	if w.Condition == "" {
		t.Error("expected non-empty condition")
	}
}

func TestGetWeather_Deterministic(t *testing.T) {
	w1, _ := GetWeather("Paris")
	w2, _ := GetWeather("Paris")
	if w1.Condition != w2.Condition {
		t.Errorf("expected deterministic result, got %s and %s", w1.Condition, w2.Condition)
	}
}

func TestFormatWeatherReport(t *testing.T) {
	w := &WeatherCondition{
		Condition:    "Sunny",
		TemperatureC: 25,
		Humidity:     45,
		Description:  "Nice day.",
	}
	report := FormatWeatherReport("Berlin", w)
	if !strings.Contains(report, "Berlin") {
		t.Error("expected report to contain location")
	}
	if !strings.Contains(report, "Sunny") {
		t.Error("expected report to contain condition")
	}
	if !strings.Contains(report, "25") {
		t.Error("expected report to contain temperature")
	}
}

func TestIsWeatherCommand(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"!weather Tokyo", true},
		{"!weather", true},
		{"!poll something", false},
		{"weather Tokyo", false},
	}
	for _, c := range cases {
		if got := isWeatherCommand(c.input); got != c.expected {
			t.Errorf("isWeatherCommand(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestHandleWeatherCommand_NoLocation(t *testing.T) {
	resp := handleWeatherCommand("!weather")
	if !strings.Contains(resp, "Usage") {
		t.Errorf("expected usage hint, got: %s", resp)
	}
}

func TestHandleWeatherCommand_WithLocation(t *testing.T) {
	resp := handleWeatherCommand("!weather Sydney")
	if !strings.Contains(resp, "Sydney") {
		t.Errorf("expected location in response, got: %s", resp)
	}
}
