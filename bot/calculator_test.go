package bot

import (
	"strings"
	"testing"
)

func TestEvaluateSimple_Addition(t *testing.T) {
	r, err := EvaluateSimple("3 + 4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 7 {
		t.Errorf("expected 7, got %v", r.Result)
	}
}

func TestEvaluateSimple_Subtraction(t *testing.T) {
	r, err := EvaluateSimple("10 - 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 7 {
		t.Errorf("expected 7, got %v", r.Result)
	}
}

func TestEvaluateSimple_Multiplication(t *testing.T) {
	r, err := EvaluateSimple("6 * 7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 42 {
		t.Errorf("expected 42, got %v", r.Result)
	}
}

func TestEvaluateSimple_Division(t *testing.T) {
	r, err := EvaluateSimple("15 / 4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 3.75 {
		t.Errorf("expected 3.75, got %v", r.Result)
	}
}

func TestEvaluateSimple_DivisionByZero(t *testing.T) {
	_, err := EvaluateSimple("5 / 0")
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestEvaluateSimple_Power(t *testing.T) {
	r, err := EvaluateSimple("2 ^ 8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 256 {
		t.Errorf("expected 256, got %v", r.Result)
	}
}

func TestEvaluateSimple_Empty(t *testing.T) {
	_, err := EvaluateSimple("")
	if err == nil {
		t.Error("expected error for empty expression")
	}
}

func TestEvaluateSimple_Invalid(t *testing.T) {
	_, err := EvaluateSimple("abc + 1")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestIsCalcCommand(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"!calc 3 + 4", true},
		{"!calculate 2 * 2", true},
		{"!calc", true},
		{"!hello", false},
		{"calc 1 + 1", false},
	}
	for _, c := range cases {
		got := isCalcCommand(c.input)
		if got != c.want {
			t.Errorf("isCalcCommand(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestHandleCalcCommand_NoArgs(t *testing.T) {
	reply := handleCalcCommand("!calc")
	if !strings.Contains(reply, "Usage") {
		t.Errorf("expected usage hint, got: %s", reply)
	}
}

func TestHandleCalcCommand_Valid(t *testing.T) {
	reply := handleCalcCommand("!calc 3 + 4")
	if !strings.Contains(reply, "7") {
		t.Errorf("expected result 7 in reply, got: %s", reply)
	}
}

func TestFormatCalcResult_Integer(t *testing.T) {
	r, _ := EvaluateSimple("6 * 7")
	out := FormatCalcResult(r)
	if !strings.Contains(out, "42") {
		t.Errorf("expected 42 in output, got: %s", out)
	}
}
