package bot

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
)

// Joke represents a single joke entry.
type Joke struct {
	ID   int
	Text string
}

// JokeStore manages a collection of jokes.
type JokeStore struct {
	mu     sync.RWMutex
	jokes  []Joke
	nextID int
}

var defaultJokes = []string{
	"Why do programmers prefer dark mode? Because light attracts bugs!",
	"Why did the developer go broke? Because he used up all his cache.",
	"A SQL query walks into a bar, walks up to two tables and asks... 'Can I join you?'",
	"How many programmers does it take to change a light bulb? None, that's a hardware problem.",
	"Why do Java developers wear glasses? Because they don't C#.",
}

// NewJokeStore creates a JokeStore pre-loaded with default jokes.
func NewJokeStore() *JokeStore {
	s := &JokeStore{}
	for _, text := range defaultJokes {
		s.jokes = append(s.jokes, Joke{ID: s.nextID, Text: text})
		s.nextID++
	}
	return s
}

// Add adds a new joke and returns its ID.
func (s *JokeStore) Add(text string) (int, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, errors.New("joke text cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.jokes = append(s.jokes, Joke{ID: id, Text: text})
	s.nextID++
	return id, nil
}

// Random returns a random joke or an error if the store is empty.
func (s *JokeStore) Random() (Joke, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.jokes) == 0 {
		return Joke{}, errors.New("no jokes available")
	}
	return s.jokes[rand.Intn(len(s.jokes))], nil
}

// List returns all jokes.
func (s *JokeStore) List() []Joke {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Joke, len(s.jokes))
	copy(out, s.jokes)
	return out
}

// Remove deletes a joke by ID. Returns an error if not found.
func (s *JokeStore) Remove(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, j := range s.jokes {
		if j.ID == id {
			s.jokes = append(s.jokes[:i], s.jokes[i+1:]...)
			return nil
		}
	}
	return errors.New("joke not found")
}
