package bot

import (
	"fmt"
	"sync"
	"time"
)

// Reminder represents a scheduled reminder
type Reminder struct {
	ID        string
	UserID    string
	ChannelID string
	Message   string
	TriggerAt time.Time
	Timer     *time.Timer
}

// ReminderStore manages active reminders
type ReminderStore struct {
	mu      sync.Mutex
	items   map[string]*Reminder
	nextID  int
}

// NewReminderStore creates a new ReminderStore
func NewReminderStore() *ReminderStore {
	return &ReminderStore{
		items: make(map[string]*Reminder),
	}
}

// Add schedules a new reminder and returns its ID
func (rs *ReminderStore) Add(userID, channelID, message string, duration time.Duration, onTrigger func(*Reminder)) string {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.nextID++
	id := fmt.Sprintf("%s-%d", userID, rs.nextID)

	r := &Reminder{
		ID:        id,
		UserID:    userID,
		ChannelID: channelID,
		Message:   message,
		TriggerAt: time.Now().Add(duration),
	}

	r.Timer = time.AfterFunc(duration, func() {
		onTrigger(r)
		rs.Remove(id)
	})

	rs.items[id] = r
	return id
}

// Remove cancels and deletes a reminder by ID
func (rs *ReminderStore) Remove(id string) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	r, ok := rs.items[id]
	if !ok {
		return false
	}
	r.Timer.Stop()
	delete(rs.items, id)
	return true
}

// List returns all reminders for a given user
func (rs *ReminderStore) List(userID string) []*Reminder {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var result []*Reminder
	for _, r := range rs.items {
		if r.UserID == userID {
			result = append(result, r)
		}
	}
	return result
}
