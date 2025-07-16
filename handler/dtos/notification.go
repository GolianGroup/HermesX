package dtos

import "github.com/gocql/gocql"

type Notification struct {
	ProfileId gocql.UUID `json:"profile_id" validate:"required,uuid"`
	Title     string     `json:"title" validate:"required,min=3"`
	Message   string     `json:"message" validate:"required,min=3"`
	Event     string     `json:"event" validate:"required"`
}

type ReadNotification struct {
	ProfileId      gocql.UUID `json:"profile_id" validate:"required,uuid"`
	NotificationId gocql.UUID `json:"notification_id" validate:"required"`
}

type ReadBroadcast struct {
	ProfileId   gocql.UUID `json:"profile_id" validate:"required,uuid"`
	BroadcastId gocql.UUID `json:"broadcast_id" validate:"required,uuid"`
}

type Broadcast struct {
	Title   string `json:"title" validate:"required,min=3"`
	Message string `json:"message" validate:"required,min=3"`
	Event   string `json:"event" validate:"required,min=3"`
}
