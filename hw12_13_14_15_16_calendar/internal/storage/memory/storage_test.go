package memorystorage

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		operation int
		input     storage.Event
		expected  []storage.Event
	}{
		{operation: 0, input: storage.Event{
			ID:             1,
			Title:          "test api requests in unit test",
			DataStart:      time.Now(),
			DataEnd:        time.Now().Add(60 * time.Minute),
			Description:    "test api requests in unit test: create, delete, update, read.",
			IdUser:         1,
			TimeUntilEvent: 15,
		}},
		{operation: 1, expected: []storage.Event{{
			ID:             1,
			Title:          "test api requests in unit test",
			DataStart:      time.Now(),
			DataEnd:        time.Now().Add(60 * time.Minute),
			Description:    "test api requests in unit test: create, delete, update, read.",
			IdUser:         1,
			TimeUntilEvent: 15,
		}}},
		{operation: 2, input: storage.Event{
			ID:             1,
			Title:          "test api requests in unit test",
			DataStart:      time.Now().Add(60 * time.Minute),
			DataEnd:        time.Now().Add(120 * time.Minute),
			Description:    "test api requests in unit test: create, delete, update, read.",
			IdUser:         1,
			TimeUntilEvent: 20,
		}},
		{operation: 1, expected: []storage.Event{{
			ID:             1,
			Title:          "test api requests in unit test",
			DataStart:      time.Now().Add(60 * time.Minute),
			DataEnd:        time.Now().Add(120 * time.Minute),
			Description:    "test api requests in unit test: create, delete, update, read.",
			IdUser:         1,
			TimeUntilEvent: 20,
		}}},
		{operation: 3, input: storage.Event{ID: 1}},
		{operation: 1, expected: []storage.Event{}},
	}
	s := New()
	for _, test := range tests {
		fmt.Println("xyi")
		switch test.operation {
		case 0:
			err := s.CreateEvent(nil, test.input)
			require.NoError(t, err)
		case 1:
			events, err := s.GetEvents()
			require.NoError(t, err)
			require.Equal(t, len(test.expected), len(events))
		case 2:
			id := test.input.ID
			err := s.UpdateEvent(id, test.input)
			require.NoError(t, err)
		case 3:
			id := test.input.ID
			err := s.DeleteEvent(id)
			require.NoError(t, err)
		}
	}
}
func TestStorageError(t *testing.T) {
	s := New()
	tests := []struct {
		operation int
		err       error
	}{
		{operation: 1, err: errors.New("there is no such id.: 1")},
		{operation: 2, err: errors.New("the parameter is not specified correctly: the title is empty")},
		{operation: 3, err: errors.New("there is no such id.: 1")},
	}
	for _, test := range tests {
		switch test.operation {
		case 1:
			err := s.UpdateEvent(1, storage.Event{})
			require.EqualError(t, err, test.err.Error())
		case 2:
			err := s.CreateEvent(nil, storage.Event{})
			err = s.UpdateEvent(1, storage.Event{})
			require.EqualError(t, err, test.err.Error())
		case 3:
			_ = s.DeleteEvent(1)
			err := s.DeleteEvent(1)
			require.EqualError(t, err, test.err.Error())
		}
	}
}
