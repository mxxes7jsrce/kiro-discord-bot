package bot

import (
	"sync"
	"testing"
	"time"
)

func TestReminderStore_AddAndTrigger(t *testing.T) {
	store := NewReminderStore()
	var wg sync.WaitGroup
	wg.Add(1)

	triggered := false
	store.Add("user1", "chan1", "hello", 50*time.Millisecond, func(r *Reminder) {
		triggered = true
		if r.Message != "hello" {
			t.Errorf("expected message 'hello', got '%s'", r.Message)
		}
		wg.Done()
	})

	wg.Wait()
	if !triggered {
		t.Error("reminder was not triggered")
	}
}

func TestReminderStore_Remove(t *testing.T) {
	store := NewReminderStore()
	triggered := false

	id := store.Add("user1", "chan1", "cancel me", 200*time.Millisecond, func(r *Reminder) {
		triggered = true
	})

	removed := store.Remove(id)
	if !removed {
		t.Error("expected Remove to return true")
	}

	time.Sleep(300 * time.Millisecond)
	if triggered {
		t.Error("reminder should not have triggered after cancellation")
	}
}

func TestReminderStore_List(t *testing.T) {
	store := NewReminderStore()

	store.Add("userA", "chan1", "msg1", 1*time.Hour, func(r *Reminder) {})
	store.Add("userA", "chan2", "msg2", 1*time.Hour, func(r *Reminder) {})
	store.Add("userB", "chan1", "msg3", 1*time.Hour, func(r *Reminder) {})

	listA := store.List("userA")
	if len(listA) != 2 {
		t.Errorf("expected 2 reminders for userA, got %d", len(listA))
	}

	listB := store.List("userB")
	if len(listB) != 1 {
		t.Errorf("expected 1 reminder for userB, got %d", len(listB))
	}

	// cleanup
	for _, r := range listA {
		store.Remove(r.ID)
	}
	for _, r := range listB {
		store.Remove(r.ID)
	}
}

func TestReminderStore_RemoveNonExistent(t *testing.T) {
	store := NewReminderStore()
	if store.Remove("nonexistent") {
		t.Error("expected Remove to return false for nonexistent ID")
	}
}
