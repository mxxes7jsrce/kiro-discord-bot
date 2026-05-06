package bot

import (
	"strings"
	"testing"
)

func TestParseRPSChoice_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected RPSChoice
	}{
		{"rock", Rock},
		{"Rock", Rock},
		{"r", Rock},
		{"paper", Paper},
		{"P", Paper},
		{"scissors", Scissors},
		{"S", Scissors},
	}
	for _, tc := range cases {
		c, ok := parseRPSChoice(tc.input)
		if !ok {
			t.Errorf("parseRPSChoice(%q) returned not ok", tc.input)
		}
		if c != tc.expected {
			t.Errorf("parseRPSChoice(%q) = %v, want %v", tc.input, c, tc.expected)
		}
	}
}

func TestParseRPSChoice_Invalid(t *testing.T) {
	invalids := []string{"", "lizard", "spock", "123"}
	for _, input := range invalids {
		_, ok := parseRPSChoice(input)
		if ok {
			t.Errorf("parseRPSChoice(%q) should be invalid", input)
		}
	}
}

func TestDetermineRPSOutcome(t *testing.T) {
	cases := []struct {
		player, bot RPSChoice
		expected    string
	}{
		{Rock, Rock, "tie"},
		{Paper, Paper, "tie"},
		{Scissors, Scissors, "tie"},
		{Rock, Scissors, "win"},
		{Paper, Rock, "win"},
		{Scissors, Paper, "win"},
		{Rock, Paper, "lose"},
		{Paper, Scissors, "lose"},
		{Scissors, Rock, "lose"},
	}
	for _, tc := range cases {
		got := determineRPSOutcome(tc.player, tc.bot)
		if got != tc.expected {
			t.Errorf("determineRPSOutcome(%v, %v) = %q, want %q", tc.player, tc.bot, got, tc.expected)
		}
	}
}

func TestPlayRPS_EmptyInput(t *testing.T) {
	_, err := PlayRPS("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestPlayRPS_InvalidInput(t *testing.T) {
	_, err := PlayRPS("lizard")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestPlayRPS_ValidOutcome(t *testing.T) {
	for i := 0; i < 20; i++ {
		result, err := PlayRPS("rock")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Outcome != "win" && result.Outcome != "lose" && result.Outcome != "tie" {
			t.Errorf("unexpected outcome: %q", result.Outcome)
		}
	}
}

func TestIsRPSCommand(t *testing.T) {
	if !isRPSCommand([]string{"rps", "rock"}) {
		t.Error("expected rps to be recognized")
	}
	if isRPSCommand([]string{"dice", "rock"}) {
		t.Error("dice should not be rps command")
	}
	if isRPSCommand([]string{}) {
		t.Error("empty args should return false")
	}
}

func TestDispatchRPSCommand_NoArgs(t *testing.T) {
	out := dispatchRPSCommand([]string{"rps"})
	if !strings.Contains(out, "Usage") {
		t.Errorf("expected usage hint, got %q", out)
	}
}

func TestDispatchRPSCommand_InvalidChoice(t *testing.T) {
	out := dispatchRPSCommand([]string{"rps", "banana"})
	if !strings.Contains(out, "❌") {
		t.Errorf("expected error emoji, got %q", out)
	}
}

func TestDispatchRPSCommand_ValidChoice(t *testing.T) {
	out := dispatchRPSCommand([]string{"rps", "paper"})
	if out == "" {
		t.Error("expected non-empty response")
	}
}
