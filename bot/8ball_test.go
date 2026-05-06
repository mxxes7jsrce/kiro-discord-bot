package bot

import (
	"strings"
	"testing"
)

func TestAskEightBall_EmptyQuestion(t *testing.T) {
	result := AskEightBall("")
	if result != "" {
		t.Errorf("expected empty string for empty question, got %q", result)
	}
}

func TestAskEightBall_NonEmptyQuestion(t *testing.T) {
	question := "Will this test pass?"
	result := AskEightBall(question)
	if result == "" {
		t.Error("expected a non-empty response for a valid question")
	}
	// Verify the response is one of the known responses.
	found := false
	for _, r := range eightBallResponses {
		if r == result {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("response %q not in known responses", result)
	}
}

func TestIsEightBallCommand(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"8ball", true},
		{"8b", true},
		{"ball", false},
		{"dice", false},
		{"", false},
	}
	for _, c := range cases {
		result := isEightBallCommand(c.input)
		if result != c.expected {
			t.Errorf("isEightBallCommand(%q) = %v, want %v", c.input, result, c.expected)
		}
	}
}

func TestHandleEightBall_NoQuestion(t *testing.T) {
	result := handleEightBall("")
	if !strings.Contains(result, "Please ask a question") {
		t.Errorf("expected help message for empty question, got %q", result)
	}
}

func TestHandleEightBall_WithQuestion(t *testing.T) {
	result := handleEightBall("Is Go awesome?")
	if !strings.HasPrefix(result, "🎱 **") {
		t.Errorf("expected formatted response, got %q", result)
	}
	if !strings.HasSuffix(result, "**") {
		t.Errorf("expected response to end with '**', got %q", result)
	}
}
