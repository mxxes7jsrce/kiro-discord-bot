package bot

import (
	"strings"
	"testing"
)

func TestDispatchJokeCommand_Random(t *testing.T) {
	s := NewJokeStore()
	out := dispatchJokeCommand(s, "!joke")
	if !strings.Contains(out, "Joke #") {
		t.Errorf("expected joke output, got: %s", out)
	}
}

func TestDispatchJokeCommand_Add(t *testing.T) {
	s := &JokeStore{}
	out := dispatchJokeCommand(s, "!joke add Why did the cat sit on the computer?")
	if !strings.Contains(out, "added") {
		t.Errorf("expected confirmation, got: %s", out)
	}
	if len(s.List()) != 1 {
		t.Errorf("expected 1 joke in store, got %d", len(s.List()))
	}
}

func TestDispatchJokeCommand_AddMissingArgs(t *testing.T) {
	s := &JokeStore{}
	out := dispatchJokeCommand(s, "!joke add")
	if !strings.Contains(out, "Usage") {
		t.Errorf("expected usage hint, got: %s", out)
	}
}

func TestDispatchJokeCommand_Remove(t *testing.T) {
	s := &JokeStore{}
	id, _ := s.Add("Temp joke")
	out := dispatchJokeCommand(s, strings.Join([]string{"!joke remove", strings.TrimSpace(strings.Join([]string{""}, "") + fmt.Sprintf("%d", id))}, " "))
	_ = out // just ensure no panic
}

func TestDispatchJokeCommand_List(t *testing.T) {
	s := &JokeStore{}
	s.Add("Joke one")
	s.Add("Joke two")
	out := dispatchJokeCommand(s, "!joke list")
	if !strings.Contains(out, "Joke one") || !strings.Contains(out, "Joke two") {
		t.Errorf("expected both jokes listed, got: %s", out)
	}
}

func TestDispatchJokeCommand_UnknownSubCommand(t *testing.T) {
	s := NewJokeStore()
	out := dispatchJokeCommand(s, "!joke unknown")
	if !strings.Contains(out, "Joke Commands") {
		t.Errorf("expected help fallback, got: %s", out)
	}
}
