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
		return errors.New("Заголовок пустой ")
	}
	if e.DataEnd.IsZero() {
		return errors.New("Время окончание события не указан ")
	}
	if e.IdUser == 0 {
		return errors.New("Не указан id пользователя ")
	}
	if !e.DataEnd.After(e.DataStart) {
		return errors.New("Не правильно указано время ")
	}
	return nil
}
