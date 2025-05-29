package app

import (
	"context"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
}

type Storage interface {
	CreateEvent(event storage.Event) error
	DeleteEvent(id int) error
	GetEvents() ([]storage.Event, error)
	UpdateEvent(id int, event storage.Event) error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {

	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
