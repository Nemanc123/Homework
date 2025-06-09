package memorystorage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu   sync.RWMutex //nolint:unused
	id   int
	list map[int]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:   sync.RWMutex{},
		id:   0,
		list: make(map[int]storage.Event),
	}
}
func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	event.ID = strconv.Itoa(s.id)
	if event.DataStart.IsZero() {
		event.DataStart = time.Now()
	}
	s.list[s.id] = event
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.list[id]; !ok {
		return errors.New(fmt.Sprintf("there is no such id.: %d", id))
	}
	delete(s.list, id)
	return nil
}

func (s *Storage) GetEvents() ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var events = make([]storage.Event, 0, len(s.list))
	for _, k := range s.list {
		events = append(events, k)
	}
	return events, nil
}

func (s *Storage) UpdateEvent(id int, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.list[id]; !ok {
		return errors.New(fmt.Sprintf("there is no such id.: %d", id))
	}
	if err := event.Validate(); err != nil {
		return fmt.Errorf("the parameter is not specified correctly: %w", err)
	}
	s.list[s.id] = event
	return nil
}
