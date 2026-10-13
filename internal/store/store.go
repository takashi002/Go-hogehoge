package store

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Mood string

const (
	MoodCalm   Mood = "calm"
	MoodStirred Mood = "stirred"
	MoodStorm  Mood = "storm"
)

type Entry struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	Mood      Mood      `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	mu      sync.RWMutex
	entries []Entry
}

func New() *Store {
	s := &Store{entries: make([]Entry, 0)}
	s.seed()
	return s
}

func (s *Store) seed() {
	now := time.Now()
	s.entries = []Entry{
		{
			ID:        uuid.NewString(),
			Body:      "窓を開けたら、遠い潮の匂いがした。今日は小さく始める。",
			Mood:      MoodCalm,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        uuid.NewString(),
			Body:      "締め切りの波が寄せてきた。呼吸を整えて、ひとつだけ片づける。",
			Mood:      MoodStirred,
			CreatedAt: now.Add(-26 * time.Hour),
		},
	}
}

func (s *Store) List() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *Store) Add(body string, mood Mood) Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := Entry{
		ID:        uuid.NewString(),
		Body:      body,
		Mood:      mood,
		CreatedAt: time.Now(),
	}
	s.entries = append([]Entry{e}, s.entries...)
	return e
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.entries {
		if e.ID == id {
			s.entries = append(s.entries[:i], s.entries[i+1:]...)
			return true
		}
	}
	return false
}

func ValidMood(m string) (Mood, bool) {
	switch Mood(m) {
	case MoodCalm, MoodStirred, MoodStorm:
		return Mood(m), true
	default:
		return "", false
	}
}
