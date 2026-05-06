package bot

import (
	"strings"
	"testing"
)

func TestDispatchCalcCommand_Addition(t *testing.T) {
	reply := dispatchCalcCommand("!calc 10 + 5")
	if !strings.Contains(reply, "15") {
		t.Errorf("expected 15 in reply, got: %s", reply)
	}
}

func TestDispatchCalcCommand_Division(t *testing.T) {
	reply := dispatchCalcCommand("!calc 9 / 3")
	if !strings.Contains(reply, "3") {
		t.Errorf("expected 3 in reply, got: %s", reply)
	}
}

func TestDispatchCalcCommand_DivisionByZero(t *testing.T) {
	reply := dispatchCalcCommand("!calc 1 / 0")
	if !strings.Contains(reply, "❌") {
		t.Errorf("expected error emoji in reply, got: %s", reply)
	}
}

func TestDispatchCalcCommand_MissingArgs(t *testing.T) {
	reply := dispatchCalcCommand("!calculate")
	if !strings.Contains(strings.ToLower(reply), "usage") {
		t.Errorf("expected usage hint, got: %s", reply)
	}
}

func TestDispatchCalcCommand_Power(t *testing.T) {
	reply := dispatchCalcCommand("!calc 2 ^ 10")
	if !strings.Contains(reply, "1024") {
		t.Errorf("expected 1024 in reply, got: %s", reply)
	}
}

func TestDispatchCalcCommand_Float(t *testing.T) {
	reply := dispatchCalcCommand("!calc 1 / 3")
	if !strings.Contains(reply, "0.333") {
		t.Errorf("expected fractional result in reply, got: %s", reply)
	}
}
