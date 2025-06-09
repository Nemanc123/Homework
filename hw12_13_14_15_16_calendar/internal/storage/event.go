package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID             string
	Title          string
	DataStart      time.Time
	DataEnd        time.Time
	Description    string
	IdUser         int
	TimeUntilEvent int
}

func (e *Event) Validate() error {
	if e.Title == "" {
		return errors.New("the title is empty")
	}
	if e.DataEnd.IsZero() {
		return errors.New("the end time of the event is not specified")
	}
	if e.IdUser == 0 {
		return errors.New("user id is not specified")
	}
	if !e.DataEnd.After(e.DataStart) {
		return errors.New("the time is not specified correctly")
	}
	return nil
}
