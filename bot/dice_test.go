package bot

import (
	"strings"
	"testing"
)

func TestParseDiceNotation_Valid(t *testing.T) {
	cases := []struct {
		input      string
		wantDice   int
		wantSides  int
	}{
		{"2d6", 2, 6},
		{"1d20", 1, 20},
		{"d10", 1, 10},
		{"3D8", 3, 8},
	}
	for _, c := range cases {
		d, s, err := ParseDiceNotation(c.input)
		if err != nil {
			t.Errorf("ParseDiceNotation(%q) unexpected error: %v", c.input, err)
			continue
		}
		if d != c.wantDice || s != c.wantSides {
			t.Errorf("ParseDiceNotation(%q) = (%d, %d), want (%d, %d)", c.input, d, s, c.wantDice, c.wantSides)
		}
	}
}

func TestParseDiceNotation_Invalid(t *testing.T) {
	invalid := []string{"abc", "2x6", "d", "2d"}
	for _, input := range invalid {
		_, _, err := ParseDiceNotation(input)
		if err == nil {
			t.Errorf("ParseDiceNotation(%q) expected error, got nil", input)
		}
	}
}

func TestRollDice_RangeAndTotal(t *testing.T) {
	result, err := RollDice(4, 6)
	if err != nil {
		t.Fatalf("RollDice(4,6) unexpected error: %v", err)
	}
	if len(result.Rolls) != 4 {
		t.Errorf("expected 4 rolls, got %d", len(result.Rolls))
	}
	sum := 0
	for _, r := range result.Rolls {
		if r < 1 || r > 6 {
			t.Errorf("roll value %d out of range [1,6]", r)
		}
		sum += r
	}
	if sum != result.Total {
		t.Errorf("total mismatch: sum=%d, Total=%d", sum, result.Total)
	}
}

func TestRollDice_InvalidArgs(t *testing.T) {
	if _, err := RollDice(0, 6); err == nil {
		t.Error("expected error for dice=0")
	}
	if _, err := RollDice(1, 1); err == nil {
		t.Error("expected error for sides=1")
	}
	if _, err := RollDice(21, 6); err == nil {
		t.Error("expected error for dice=21")
	}
}

func TestHandleRollDice_Default(t *testing.T) {
	out := handleRollDice([]string{})
	if !strings.Contains(out, "d6") {
		t.Errorf("expected default d6 roll, got: %s", out)
	}
}

func TestHandleRollDice_Custom(t *testing.T) {
	out := handleRollDice([]string{"2d10"})
	if !strings.Contains(out, "2d10") {
		t.Errorf("expected 2d10 in output, got: %s", out)
	}
}

func TestIsDiceCommand(t *testing.T) {
	if !isDiceCommand("!roll 2d6") {
		t.Error("expected !roll to be a dice command")
	}
	if isDiceCommand("!help") {
		t.Error("expected !help to not be a dice command")
	}
}
