package bot

import (
	"strings"
	"testing"
)

func TestGetCompliment_NoTarget(t *testing.T) {
	c := GetCompliment("")
	if c == "" {
		t.Error("expected a non-empty compliment")
	}
	found := false
	for _, dc := range defaultCompliments {
		if c == dc {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("compliment %q not in defaultCompliments", c)
	}
}

func TestGetCompliment_WithTarget(t *testing.T) {
	c := GetCompliment("Alice")
	if !strings.HasPrefix(c, "Alice, ") {
		t.Errorf("expected compliment to start with 'Alice, ', got %q", c)
	}
}

func TestIsComplimentCommand(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"!compliment", true},
		{"!compliment @Bob", true},
		{"!help", false},
		{"compliment", false},
		{"", false},
	}
	for _, tc := range cases {
		got := isComplimentCommand(tc.input)
		if got != tc.want {
			t.Errorf("isComplimentCommand(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestHandleComplimentCommand_NoArgs(t *testing.T) {
	res := handleComplimentCommand([]string{})
	if res == "" {
		t.Error("expected non-empty response")
	}
}

func TestHandleComplimentCommand_WithTarget(t *testing.T) {
	res := handleComplimentCommand([]string{"@Charlie"})
	if !strings.HasPrefix(res, "@Charlie, ") {
		t.Errorf("expected response to start with '@Charlie, ', got %q", res)
	}
}

func TestDispatchComplimentCommand(t *testing.T) {
	res := dispatchComplimentCommand("!compliment @Dave")
	if !strings.HasPrefix(res, "@Dave, ") {
		t.Errorf("expected response to start with '@Dave, ', got %q", res)
	}
}

func TestDispatchComplimentCommand_NoTarget(t *testing.T) {
	res := dispatchComplimentCommand("!compliment")
	if res == "" {
		t.Error("expected non-empty response for no-target dispatch")
	}
}
