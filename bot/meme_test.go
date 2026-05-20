package bot

import (
	"strings"
	"testing"
)

func TestGetMemeTemplateNames(t *testing.T) {
	names := GetMemeTemplateNames()
	if len(names) == 0 {
		t.Fatal("expected at least one meme template")
	}
	for _, n := range names {
		if n == "" {
			t.Error("template name should not be empty")
		}
	}
}

func TestGenerateMeme_KnownTemplate(t *testing.T) {
	result, err := GenerateMeme("onedoesnot", []string{"write tests"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "write tests") {
		t.Errorf("expected result to contain arg, got: %s", result)
	}
}

func TestGenerateMeme_UnknownTemplate(t *testing.T) {
	_, err := GenerateMeme("notreal", []string{"arg"})
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
}

func TestGenerateMeme_TooFewArgs(t *testing.T) {
	_, err := GenerateMeme("drake", []string{"only one arg"})
	if err == nil {
		t.Fatal("expected error when too few args provided")
	}
}

func TestGenerateMeme_Random(t *testing.T) {
	result, err := GenerateMeme("random", []string{"foo", "bar", "baz"})
	if err != nil {
		t.Fatalf("unexpected error for random meme: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result for random meme")
	}
}

func TestIsMemeCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!meme", true},
		{"!meme list", true},
		{"!meme drake yes no", true},
		{"!poll something", false},
		{"meme", false},
		{"!MEME list", true},
	}
	for _, tc := range tests {
		got := isMemeCommand(tc.input)
		if got != tc.expected {
			t.Errorf("isMemeCommand(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestDispatchMemeCommand_List(t *testing.T) {
	result := dispatchMemeCommand("!meme list")
	if !strings.Contains(result, "drake") {
		t.Errorf("expected list to contain 'drake', got: %s", result)
	}
}

func TestDispatchMemeCommand_UnknownTemplate(t *testing.T) {
	result := dispatchMemeCommand("!meme faketemplate arg1")
	if !strings.Contains(result, "❌") {
		t.Errorf("expected error emoji in response, got: %s", result)
	}
}

func TestDispatchMemeCommand_NoArgs(t *testing.T) {
	result := dispatchMemeCommand("!meme")
	if !strings.Contains(result, "🖼️") {
		t.Errorf("expected meme emoji in response, got: %s", result)
	}
}
