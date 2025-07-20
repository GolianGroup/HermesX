package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"hermesx/internal/ent"
	"hermesx/internal/ent/event"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type EventRepository interface {
	Create(ctx context.Context, eventData *ent.Event) error
	GetByName(ctx context.Context, name string) (*ent.Event, error)
	GetAll(ctx context.Context) ([]*ent.Event, error)
	DoesExist(ctx context.Context, name string) (bool, error)
	Remove(ctx context.Context, name string) error
}

type eventRepository struct {
	db *ent.Client
}

func NewEventRepository(db *ent.Client) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, eventData *ent.Event) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	creator := r.db.Event.
		Create().
		SetID(eventData.ID).
		SetName(eventData.Name).
		SetIsCritical(eventData.IsCritical).
		SetIsProtected(eventData.IsProtected).
		SetPreferredChannel(eventData.PreferredChannel).
		SetTemplate(eventData.Template)

	// Set metadata only if it's provided
	if eventData.Metadata != nil {
		creator.SetMetadata(eventData.Metadata)
	} else {
		// Set empty map as default value
		creator.SetMetadata(map[string]any{})
	}

	_, err := creator.Save(ctx)
	if err != nil {
		log.Println("Error creating event:", err)
		return ErrDatabase
	}
	return nil
}

func (r *eventRepository) GetByName(ctx context.Context, name string) (*ent.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	eventData, err := r.db.Event.Query().Where(event.Name(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrEventNotFound
		}
		return nil, ErrDatabase
	}

	return eventData, nil
}

func (r *eventRepository) DoesExist(ctx context.Context, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	count, err := r.db.Event.Query().Where(event.Name(name)).Count(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, ErrDatabase
	}
	return count > 0, nil
}

func (r *eventRepository) GetAll(ctx context.Context) ([]*ent.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	events, err := r.db.Event.Query().All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return []*ent.Event{}, nil
		}
		return nil, ErrDatabase
	}

	return events, nil
}

func (r *eventRepository) Remove(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.Event.Delete().Where(event.Name(name)).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrEventNotFound
		}
		return ErrDatabase
	}

	return nil
}
