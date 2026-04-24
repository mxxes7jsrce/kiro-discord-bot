package bot

import (
	"math/rand"
	"sync"
	"time"
)

// TriviaQuestion represents a single trivia question with options and answer.
type TriviaQuestion struct {
	Question string
	Options  []string
	Answer   int // 0-based index into Options
}

// TriviaSession holds an active trivia game in a channel.
type TriviaSession struct {
	Question  *TriviaQuestion
	ChannelID string
	StartedAt time.Time
	Answered  bool
}

// TriviaStore manages active trivia sessions per channel.
type TriviaStore struct {
	mu       sync.Mutex
	sessions map[string]*TriviaSession
	questions []*TriviaQuestion
}

// NewTriviaStore creates a TriviaStore pre-loaded with sample questions.
func NewTriviaStore() *TriviaStore {
	return &TriviaStore{
		sessions: make(map[string]*TriviaSession),
		questions: defaultQuestions(),
	}
}

// Start begins a new trivia session in the given channel.
// Returns the question and true, or nil and false if a session is already active.
func (ts *TriviaStore) Start(channelID string) (*TriviaQuestion, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.sessions[channelID]; exists {
		return nil, false
	}

	q := ts.questions[rand.Intn(len(ts.questions))]
	ts.sessions[channelID] = &TriviaSession{
		Question:  q,
		ChannelID: channelID,
		StartedAt: time.Now(),
	}
	return q, true
}

// Answer checks a submitted answer index for the active session.
// Returns (correct, found). Ends the session if correct.
func (ts *TriviaStore) Answer(channelID string, idx int) (bool, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	session, exists := ts.sessions[channelID]
	if !exists {
		return false, false
	}
	if idx == session.Question.Answer {
		delete(ts.sessions, channelID)
		return true, true
	}
	return false, true
}

// Stop cancels any active trivia session in the channel.
func (ts *TriviaStore) Stop(channelID string) (*TriviaQuestion, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	session, exists := ts.sessions[channelID]
	if !exists {
		return nil, false
	}
	delete(ts.sessions, channelID)
	return session.Question, true
}

func defaultQuestions() []*TriviaQuestion {
	return []*TriviaQuestion{
		{
			Question: "What is the capital of France?",
			Options:  []string{"Berlin", "Madrid", "Paris", "Rome"},
			Answer:   2,
		},
		{
			Question: "Which planet is known as the Red Planet?",
			Options:  []string{"Venus", "Mars", "Jupiter", "Saturn"},
			Answer:   1,
		},
		{
			Question: "What is 7 × 8?",
			Options:  []string{"54", "56", "58", "64"},
			Answer:   1,
		},
		{
			Question: "Who wrote 'Romeo and Juliet'?",
			Options:  []string{"Charles Dickens", "Mark Twain", "William Shakespeare", "Jane Austen"},
			Answer:   2,
		},
	}
}
