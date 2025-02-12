package app

import (
	"context"
	"fmt"

	"github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	a.logger.Debug(fmt.Sprintf("Creating event, id: %s title %s", id, title))
	event := storage.Event{
		ID:    id,
		Title: title,
	}
	if err := a.storage.AddEvent(ctx, event); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to create event, error: %s", err))
		return err
	}

	a.logger.Info(fmt.Sprintf("Event created successfully, id: %s title %s", id, title))
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListEvents(ctx)
}
