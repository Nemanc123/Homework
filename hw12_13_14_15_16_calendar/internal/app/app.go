package app

import (
	"context"
	"fmt"
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

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	err := a.storage.CreateEvent(event)
	if err != nil {
		return fmt.Errorf("couldn't create event: %w", err)
	}
	return nil
}
func (a *App) DeleteEvent(ctx context.Context, id int) error {
	err := a.storage.DeleteEvent(id)
	if err != nil {
		return fmt.Errorf("couldn't delete event: %w", err)
	}
	return nil
}
func (a *App) GetEvent(ctx context.Context) ([]storage.Event, error) {
	resultEvents, err := a.storage.GetEvents()
	if err != nil {
		return nil, fmt.Errorf("couldn't create event: %w", err)
	}
	return resultEvents, nil
}
func (a *App) UpdateEvent(ctx context.Context, id int, event storage.Event) error {
	err := a.storage.UpdateEvent(id, event)
	if err != nil {
		return fmt.Errorf("couldn't create event: %w", err)
	}
	return nil
}

// TODO
