package bot

import (
	"fmt"
	"strings"
	"sync"
)

// PollOption represents a single option in a poll
type PollOption struct {
	Label string
	Votes map[string]bool // userID -> voted
}

// Poll represents an active poll
type Poll struct {
	Question string
	Options  []*PollOption
	Creator  string
	ChannelID string
	MessageID string
	Closed   bool
}

// PollStore manages active polls per channel
type PollStore struct {
	mu    sync.RWMutex
	polls map[string]*Poll // channelID -> Poll
}

// NewPollStore creates a new PollStore
func NewPollStore() *PollStore {
	return &PollStore{
		polls: make(map[string]*Poll),
	}
}

// Create adds a new poll for a channel, returns error if one already exists
func (ps *PollStore) Create(channelID, creator, question string, options []string) (*Poll, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.polls[channelID]; exists {
		return nil, fmt.Errorf("a poll is already active in this channel")
	}

	if len(options) < 2 || len(options) > 9 {
		return nil, fmt.Errorf("poll must have between 2 and 9 options")
	}

	pollOptions := make([]*PollOption, len(options))
	for i, label := range options {
		pollOptions[i] = &PollOption{
			Label: strings.TrimSpace(label),
			Votes: make(map[string]bool),
		}
	}

	p := &Poll{
		Question:  question,
		Options:   pollOptions,
		Creator:   creator,
		ChannelID: channelID,
	}
	ps.polls[channelID] = p
	return p, nil
}

// Vote records a vote for the given option index (1-based)
func (ps *PollStore) Vote(channelID, userID string, optionIndex int) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	p, exists := ps.polls[channelID]
	if !exists {
		return fmt.Errorf("no active poll in this channel")
	}
	if p.Closed {
		return fmt.Errorf("this poll is closed")
	}
	if optionIndex < 1 || optionIndex > len(p.Options) {
		return fmt.Errorf("invalid option, choose between 1 and %d", len(p.Options))
	}

	opt := p.Options[optionIndex-1]
	opt.Votes[userID] = true
	return nil
}

// Close ends the poll and returns results summary
func (ps *PollStore) Close(channelID string) (*Poll, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	p, exists := ps.polls[channelID]
	if !exists {
		return nil, fmt.Errorf("no active poll in this channel")
	}
	p.Closed = true
	delete(ps.polls, channelID)
	return p, nil
}

// Get returns the active poll for a channel
func (ps *PollStore) Get(channelID string) (*Poll, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	p, ok := ps.polls[channelID]
	return p, ok
}

// FormatResults returns a human-readable results string
func (p *Poll) FormatResults() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 **%s**\n", p.Question))
	total := 0
	for _, opt := range p.Options {
		total += len(opt.Votes)
	}
	for i, opt := range p.Options {
		count := len(opt.Votes)
		pct := 0
		if total > 0 {
			pct = count * 100 / total
		}
		sb.WriteString(fmt.Sprintf("%d. %s — %d vote(s) (%d%%)\n", i+1, opt.Label, count, pct))
	}
	sb.WriteString(fmt.Sprintf("Total votes: %d", total))
	return sb.String()
}
