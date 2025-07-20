package dtos

import "google.golang.org/protobuf/types/known/timestamppb"

type Event struct {
	ID               string                 `json:"id" validate:"omitempty,uuid4"`                              // ID is optional for creation, must be UUID if provided
	Name             string                 `json:"name" validate:"required,min=3,max=100"`                     // Name is required, and should be between 3 to 100 characters
	IsCritical       bool                   `json:"is_critical" default:"false"`                                // Critical is required, default is false
	IsProtected      bool                   `json:"is_protected" default:"false"`                               // Protected is required, default is false
	PreferredChannel string                 `json:"preferred_channel" validate:"required,oneof=email sms push"` // PreferredChannel is required and should be one of these values
	Template         string                 `json:"template"`                                                   // Template for email notifications
	CreatedAt        *timestamppb.Timestamp `json:"created_at,omitempty"`                                       // Optional for creation
	UpdatedAt        *timestamppb.Timestamp `json:"updated_at,omitempty"`                                       // Optional for creation
	Metadata         map[string]any         `json:"metadata,omitempty"`                                         // Optional metadata
}

type EventResponse struct {
	ID               string                 `json:"id"` // ID is required in response
	Name             string                 `json:"name"`
	IsCritical       bool                   `json:"is_critical"`
	IsProtected      bool                   `json:"is_protected"`
	PreferredChannel string                 `json:"preferred_channel"`
	Template         string                 `json:"template"`
	CreatedAt        *timestamppb.Timestamp `json:"created_at"`
	UpdatedAt        *timestamppb.Timestamp `json:"update_at"`
	Metadata         map[string]any         `json:"metadata,omitempty"`
}

type EventResponseList struct {
	Events []*EventResponse `json:"events"`
}
