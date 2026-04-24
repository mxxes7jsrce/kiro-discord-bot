package bot

import (
	"testing"
)

func TestQuoteStore_AddAndList(t *testing.T) {
	s := NewQuoteStore()
	id, err := s.Add("g1", "Hello world", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
	s.Add("g1", "Second quote", "Bob")
	list := s.List("g1")
	if len(list) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(list))
	}
	if list[0].Author != "Alice" {
		t.Errorf("expected Alice, got %s", list[0].Author)
	}
}

func TestQuoteStore_Remove(t *testing.T) {
	s := NewQuoteStore()
	s.Add("g1", "To be removed", "X")
	err := s.Remove("g1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.List("g1")) != 0 {
		t.Error("expected empty list after removal")
	}
}

func TestQuoteStore_RemoveNonExistent(t *testing.T) {
	s := NewQuoteStore()
	err := s.Remove("g1", 99)
	if err == nil {
		t.Error("expected error removing non-existent quote")
	}
}

func TestQuoteStore_Random(t *testing.T) {
	s := NewQuoteStore()
	_, err := s.Random("g1")
	if err == nil {
		t.Error("expected error on empty store")
	}
	s.Add("g1", "Only quote", "Y")
	q, err := s.Random("g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Text != "Only quote" {
		t.Errorf("unexpected quote text: %s", q.Text)
	}
}

func TestQuoteStore_EmptyText(t *testing.T) {
	s := NewQuoteStore()
	_, err := s.Add("g1", "", "Author")
	if err == nil {
		t.Error("expected error for empty text")
	}
}

func TestHandleAddQuote_WithAuthor(t *testing.T) {
	s := NewQuoteStore()
	res := handleAddQuote(s, "g1", "Great words -- Confucius")
	if res.Message != "Quote #1 added." {
		t.Errorf("unexpected message: %s", res.Message)
	}
	list := s.List("g1")
	if list[0].Author != "Confucius" {
		t.Errorf("expected Confucius, got %s", list[0].Author)
	}
}

func TestHandleAddQuote_NoAuthor(t *testing.T) {
	s := NewQuoteStore()
	res := handleAddQuote(s, "g1", "Anonymous wisdom")
	if res.Message != "Quote #1 added." {
		t.Errorf("unexpected message: %s", res.Message)
	}
	list := s.List("g1")
	if list[0].Author != "unknown" {
		t.Errorf("expected unknown, got %s", list[0].Author)
	}
}
