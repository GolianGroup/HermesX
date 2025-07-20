package presenters

import (
	"hermesx/handler/dtos"

	"github.com/gofiber/fiber/v2"
)

type Event struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	IsCritical       bool           `json:"is_critical"`
	IsProtected      bool           `json:"is_protected"`
	PreferredChannel string         `json:"preferred_channel"`
	Template         string         `json:"template"`
	Metadata         map[string]any `json:"metadata"`
}

func EventResponse(events []*dtos.EventResponse) *fiber.Map {
	resp := make([]Event, len(events))
	for i, event := range events {
		resp[i] = Event{
			ID:               event.ID,
			Name:             event.Name,
			IsCritical:       event.IsCritical,
			IsProtected:      event.IsProtected,
			PreferredChannel: event.PreferredChannel,
			Template:         event.Template,
			Metadata:         event.Metadata,
		}
	}
	return SuccessWithDataResponse(resp)
}
