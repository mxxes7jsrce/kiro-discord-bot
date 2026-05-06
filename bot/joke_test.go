package bot

import (
	"strings"
	"testing"
)

func TestJokeStore_DefaultJokes(t *testing.T) {
	s := NewJokeStore()
	if len(s.List()) == 0 {
		t.Fatal("expected default jokes to be loaded")
	}
}

func TestJokeStore_AddAndList(t *testing.T) {
	s := &JokeStore{}
	id, err := s.Add("Why did the chicken cross the road?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 0 {
		t.Errorf("expected id 0, got %d", id)
	}
	list := s.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 joke, got %d", len(list))
	}
	if list[0].Text != "Why did the chicken cross the road?" {
		t.Errorf("unexpected joke text: %s", list[0].Text)
	}
}

func TestJokeStore_AddEmptyText(t *testing.T) {
	s := &JokeStore{}
	_, err := s.Add("   ")
	if err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestJokeStore_Random(t *testing.T) {
	s := &JokeStore{}
	_, err := s.Random()
	if err == nil {
		t.Fatal("expected error when store is empty")
	}
	s.Add("Test joke")
	j, err := s.Random()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.Text != "Test joke" {
		t.Errorf("unexpected text: %s", j.Text)
	}
}

func TestJokeStore_Remove(t *testing.T) {
	s := &JokeStore{}
	id, _ := s.Add("Joke to remove")
	if err := s.Remove(id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.List()) != 0 {
		t.Error("expected empty store after removal")
	}
}

func TestJokeStore_RemoveNonExistent(t *testing.T) {
	s := &JokeStore{}
	if err := s.Remove(999); err == nil {
		t.Fatal("expected error removing non-existent joke")
	}
}

func TestHandleListJokes_Empty(t *testing.T) {
	s := &JokeStore{}
	out := handleListJokes(s)
	if !strings.Contains(out, "No jokes") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestDispatchJokeCommand_Help(t *testing.T) {
	s := NewJokeStore()
	out := dispatchJokeCommand(s, "!joke help")
	if !strings.Contains(out, "Joke Commands") {
		t.Errorf("expected help text, got: %s", out)
	}
}

func TestIsJokeCommand(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"!joke", true},
		{"!JOKE add", true},
		{"!poll", false},
		{"hello", false},
	}
	for _, c := range cases {
		if got := isJokeCommand(c.input); got != c.want {
			t.Errorf("isJokeCommand(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}
