package presenters

import (
	"hermesx/handler/dtos"

	"github.com/gofiber/fiber/v2"
)

type Event struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	IsCritical       bool           `json:"is_critical"`
	PreferredChannel string         `json:"preferred_channel"`
	Metadata         map[string]any `json:"metadata"`
}

func EventResponse(events []*dtos.EventResponse) *fiber.Map {
	resp := make([]Event, len(events))
	for i, event := range events {
		resp[i] = Event{
			ID:               event.ID,
			Name:             event.Name,
			IsCritical:       event.IsCritical,
			PreferredChannel: event.PreferredChannel,
			Metadata:         event.Metadata,
		}
	}
	return SuccessWithDataResponse(resp)
}
