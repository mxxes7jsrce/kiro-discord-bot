package bot

import (
	"testing"
)

func TestTriviaStore_StartAndAnswer(t *testing.T) {
	ts := NewTriviaStore()
	channelID := "ch-001"

	q, ok := ts.Start(channelID)
	if !ok {
		t.Fatal("expected Start to succeed")
	}
	if q == nil {
		t.Fatal("expected a question, got nil")
	}

	// Wrong answer
	wrongIdx := (q.Answer + 1) % len(q.Options)
	correct, found := ts.Answer(channelID, wrongIdx)
	if !found {
		t.Fatal("expected session to be found")
	}
	if correct {
		t.Error("expected wrong answer to return false")
	}

	// Correct answer
	correct, found = ts.Answer(channelID, q.Answer)
	if !found {
		t.Fatal("expected session to be found")
	}
	if !correct {
		t.Error("expected correct answer to return true")
	}

	// Session should be gone after correct answer
	_, found = ts.Answer(channelID, q.Answer)
	if found {
		t.Error("expected session to be removed after correct answer")
	}
}

func TestTriviaStore_DuplicateStart(t *testing.T) {
	ts := NewTriviaStore()
	channelID := "ch-002"

	_, ok := ts.Start(channelID)
	if !ok {
		t.Fatal("first start should succeed")
	}

	_, ok = ts.Start(channelID)
	if ok {
		t.Error("second start should fail when session already active")
	}
}

func TestTriviaStore_Stop(t *testing.T) {
	ts := NewTriviaStore()
	channelID := "ch-003"

	_, ok := ts.Stop(channelID)
	if ok {
		t.Error("stop on non-existent session should return false")
	}

	ts.Start(channelID)
	q, ok := ts.Stop(channelID)
	if !ok {
		t.Error("stop on active session should return true")
	}
	if q == nil {
		t.Error("stop should return the question")
	}

	// Confirm session is gone
	_, found := ts.Answer(channelID, 0)
	if found {
		t.Error("session should be removed after stop")
	}
}

func TestTriviaStore_AnswerOutOfRange(t *testing.T) {
	ts := NewTriviaStore()
	channelID := "ch-004"

	ts.Start(channelID)
	// index -1 and 99 are out of range but store only checks equality with Answer
	// The command layer validates range; store itself returns (false, true)
	_, found := ts.Answer(channelID, 99)
	if !found {
		t.Error("session should still be found after wrong answer")
	}
}
