package bot

import (
	"testing"
)

func TestPollStore_CreateAndGet(t *testing.T) {
	ps := NewPollStore()

	p, err := ps.Create("ch1", "user1", "Favourite colour?", []string{"Red", "Blue", "Green"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Question != "Favourite colour?" {
		t.Errorf("expected question 'Favourite colour?', got %q", p.Question)
	}

	got, ok := ps.Get("ch1")
	if !ok {
		t.Fatal("expected poll to exist")
	}
	if got.Question != p.Question {
		t.Errorf("expected %q, got %q", p.Question, got.Question)
	}
}

func TestPollStore_DuplicateCreate(t *testing.T) {
	ps := NewPollStore()
	_, err := ps.Create("ch1", "user1", "Q1", []string{"A", "B"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = ps.Create("ch1", "user2", "Q2", []string{"C", "D"})
	if err == nil {
		t.Error("expected error for duplicate poll, got nil")
	}
}

func TestPollStore_InvalidOptions(t *testing.T) {
	ps := NewPollStore()
	_, err := ps.Create("ch1", "user1", "Q", []string{"OnlyOne"})
	if err == nil {
		t.Error("expected error for single option poll")
	}
}

func TestPollStore_Vote(t *testing.T) {
	ps := NewPollStore()
	_, err := ps.Create("ch1", "user1", "Best lang?", []string{"Go", "Rust", "Python"})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	if err := ps.Vote("ch1", "userA", 1); err != nil {
		t.Errorf("unexpected vote error: %v", err)
	}
	if err := ps.Vote("ch1", "userB", 1); err != nil {
		t.Errorf("unexpected vote error: %v", err)
	}
	if err := ps.Vote("ch1", "userC", 2); err != nil {
		t.Errorf("unexpected vote error: %v", err)
	}

	p, _ := ps.Get("ch1")
	if len(p.Options[0].Votes) != 2 {
		t.Errorf("expected 2 votes for option 1, got %d", len(p.Options[0].Votes))
	}
	if len(p.Options[1].Votes) != 1 {
		t.Errorf("expected 1 vote for option 2, got %d", len(p.Options[1].Votes))
	}
}

func TestPollStore_VoteInvalidIndex(t *testing.T) {
	ps := NewPollStore()
	_, _ = ps.Create("ch1", "user1", "Q", []string{"A", "B"})

	if err := ps.Vote("ch1", "userA", 0); err == nil {
		t.Error("expected error for index 0")
	}
	if err := ps.Vote("ch1", "userA", 5); err == nil {
		t.Error("expected error for out-of-range index")
	}
}

func TestPollStore_Close(t *testing.T) {
	ps := NewPollStore()
	_, _ = ps.Create("ch1", "user1", "Q", []string{"Yes", "No"})
	_ = ps.Vote("ch1", "userA", 1)

	p, err := ps.Close("ch1")
	if err != nil {
		t.Fatalf("close error: %v", err)
	}
	if !p.Closed {
		t.Error("expected poll to be marked closed")
	}

	_, ok := ps.Get("ch1")
	if ok {
		t.Error("expected poll to be removed after close")
	}
}

func TestPoll_FormatResults(t *testing.T) {
	ps := NewPollStore()
	_, _ = ps.Create("ch1", "user1", "Colour?", []string{"Red", "Blue"})
	_ = ps.Vote("ch1", "u1", 1)
	_ = ps.Vote("ch1", "u2", 2)

	p, _ := ps.Get("ch1")
	result := p.FormatResults()
	if result == "" {
		t.Error("expected non-empty results string")
	}
}
