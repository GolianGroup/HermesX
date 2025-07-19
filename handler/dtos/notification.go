package dtos

import (
	"github.com/gocql/gocql"
)

type Notification struct {
	ProfileId gocql.UUID `json:"profile_id" validate:"required,uuid"`
	Title     string     `json:"title" validate:"required,min=3"`
	Message   string     `json:"message" validate:"required,min=3"`
	Event     string     `json:"event" validate:"required"`
}

type GetNotification struct {
	Id gocql.UUID `json:""`
}

type Read struct {
	ProfileId gocql.UUID `json:"profile_id" validate:"required,uuid"`
	Id        gocql.UUID `json:"id" validate:"required,uuid"`
	Type      string     `json:"type" validate:"required"`
}

type UserNotification struct {
	ProfileId gocql.UUID `params:"profileid" validate:"required,uuid"`
}

type EventByName struct {
	Name string `params:"name" validate:"required"`
}

type Broadcast struct {
	Title   string `json:"title" validate:"required,min=3"`
	Message string `json:"message" validate:"required,min=3"`
	Event   string `json:"event" validate:"required,min=3"`
}
