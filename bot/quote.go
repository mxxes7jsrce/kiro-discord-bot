package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

// Quote represents a stored quote with author info.
type Quote struct {
	ID     int
	Text   string
	Author string
}

// QuoteStore manages a collection of quotes per guild.
type QuoteStore struct {
	mu     sync.RWMutex
	quotes map[string][]Quote
	nextID map[string]int
}

// NewQuoteStore creates an empty QuoteStore.
func NewQuoteStore() *QuoteStore {
	return &QuoteStore{
		quotes: make(map[string][]Quote),
		nextID: make(map[string]int),
	}
}

// Add appends a new quote for the given guild and returns its ID.
func (s *QuoteStore) Add(guildID, text, author string) (int, error) {
	if text == "" {
		return 0, errors.New("quote text cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID[guildID]++
	id := s.nextID[guildID]
	s.quotes[guildID] = append(s.quotes[guildID], Quote{ID: id, Text: text, Author: author})
	return id, nil
}

// Remove deletes a quote by ID for the given guild.
func (s *QuoteStore) Remove(guildID string, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.quotes[guildID]
	for i, q := range list {
		if q.ID == id {
			s.quotes[guildID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("quote #%d not found", id)
}

// Random returns a random quote for the guild, or an error if none exist.
func (s *QuoteStore) Random(guildID string) (Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := s.quotes[guildID]
	if len(list) == 0 {
		return Quote{}, errors.New("no quotes stored")
	}
	return list[rand.Intn(len(list))], nil
}

// List returns all quotes for the guild.
func (s *QuoteStore) List(guildID string) []Quote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Quote, len(s.quotes[guildID]))
	copy(out, s.quotes[guildID])
	return out
}
