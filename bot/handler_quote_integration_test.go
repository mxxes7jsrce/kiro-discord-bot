package bot

import (
	"strings"
	"testing"
)

func TestIsQuoteCommand(t *testing.T) {
	tests := []struct {
		input   string
		match   bool
		expArgs string
	}{
		{"!quote add hello", true, "add hello"},
		{"!quote", true, ""},
		{"!poll create", false, ""},
		{"!quote list", true, "list"},
	}
	for _, tc := range tests {
		ok, args := isQuoteCommand(tc.input)
		if ok != tc.match {
			t.Errorf("input %q: expected match=%v, got %v", tc.input, tc.match, ok)
		}
		if ok && args != tc.expArgs {
			t.Errorf("input %q: expected args %q, got %q", tc.input, tc.expArgs, args)
		}
	}
}

func TestHandleListQuotes_Empty(t *testing.T) {
	s := NewQuoteStore()
	res := handleListQuotes(s, "g1")
	if !strings.Contains(res.Message, "No quotes") {
		t.Errorf("expected no-quotes message, got: %s", res.Message)
	}
}

func TestHandleListQuotes_WithEntries(t *testing.T) {
	s := NewQuoteStore()
	s.Add("g1", "First", "A")
	s.Add("g1", "Second", "B")
	res := handleListQuotes(s, "g1")
	if !strings.Contains(res.Message, "#1") || !strings.Contains(res.Message, "#2") {
		t.Errorf("expected both quotes listed, got: %s", res.Message)
	}
}

func TestHandleRemoveQuote_Invalid(t *testing.T) {
	s := NewQuoteStore()
	res := handleRemoveQuote(s, "g1", "notanumber")
	if !strings.Contains(res.Message, "Usage") {
		t.Errorf("expected usage message, got: %s", res.Message)
	}
}

func TestHandleRemoveQuote_NotFound(t *testing.T) {
	s := NewQuoteStore()
	res := handleRemoveQuote(s, "g1", "42")
	if !strings.Contains(res.Message, "Error") {
		t.Errorf("expected error message, got: %s", res.Message)
	}
}
