package api

import (
	"context"
	api "github.com/Calendar/hw12_13_14_15_calendar/api/gen/go"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
)

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	GetEvent(ctx context.Context) ([]storage.Event, error)
	UpdateEvent(ctx context.Context, id int, event storage.Event) error
}
type EventAPIService struct {
	storage Application
}

func NewEventAPIService(app Application) *EventAPIService {
	return &EventAPIService{app}
}
func (s *EventAPIService) Updateevent(ctx context.Context, eventId int64, event api.Event) (api.ImplResponse, error) {
	events := s.APIEventsToStorageEvents(event)
	err := s.storage.UpdateEvent(ctx, int(eventId), events)
	if err != nil {
		return api.ImplResponse{Code: 500, Body: err}, err
	}
	return api.ImplResponse{Code: 200, Body: "success"}, nil
}

func (s *EventAPIService) Deleteevent(ctx context.Context, eventId int64) (api.ImplResponse, error) {
	err := s.storage.DeleteEvent(ctx, int(eventId))
	if err != nil {
		return api.ImplResponse{Code: 500, Body: err}, err
	}
	return api.ImplResponse{Code: 200, Body: "success"}, nil
}
func (s *EventAPIService) GetAllEvents(ctx context.Context) (api.ImplResponse, error) {
	resultEvents, err := s.storage.GetEvent(ctx)
	if err != nil {
		return api.ImplResponse{Code: 500, Body: err}, err
	}
	return api.ImplResponse{Code: 200, Body: resultEvents}, nil
}
func (s *EventAPIService) Addevent(ctx context.Context, event api.Event) (api.ImplResponse, error) {
	events := s.APIEventsToStorageEvents(event)
	err := s.storage.CreateEvent(ctx, events)
	if err != nil {
		return api.ImplResponse{Code: 500, Body: err}, err
	}
	return api.ImplResponse{Code: 200, Body: "success"}, nil
}
func (s *EventAPIService) StorageEventsToAPIEvents(event storage.Event) api.Event {
	resultEventAPIStorage := api.Event{}
	resultEventAPIStorage.Id = event.ID
	resultEventAPIStorage.Title = event.Title
	resultEventAPIStorage.DataStart = event.DataStart
	resultEventAPIStorage.DataEnd = event.DataEnd
	resultEventAPIStorage.Description = event.Description
	resultEventAPIStorage.UserId = int64(event.IdUser)
	resultEventAPIStorage.TimeUntilEvent = int64(event.TimeUntilEvent)
	return resultEventAPIStorage
}

func (s *EventAPIService) APIEventsToStorageEvents(event api.Event) storage.Event {
	resultEventStorage := storage.Event{}
	resultEventStorage.ID = event.Id
	resultEventStorage.Title = event.Title
	resultEventStorage.DataStart = event.DataStart
	resultEventStorage.DataEnd = event.DataEnd
	resultEventStorage.Description = event.Description
	resultEventStorage.IdUser = int(event.UserId)
	resultEventStorage.TimeUntilEvent = int(event.TimeUntilEvent)
	return resultEventStorage
}
