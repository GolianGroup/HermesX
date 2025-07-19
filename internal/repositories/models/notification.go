package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Notification struct {
	ID        gocql.UUID `cql:"id"`
	ProfileId gocql.UUID `cql:"profile_id"`
	Title     string     `cql:"title"`
	Message   string     `cql:"message"`
}

type GetNotification struct {
	ID      gocql.UUID `cql:"id"`
	Title   string     `cql:"title"`
	Message string     `cql:"message"`
	HasRead bool       `cql:"has_read"`
}

type Broadcast struct {
	ID      gocql.UUID `cql:"id"`
	Title   string     `cql:"title"`
	Message string     `cql:"message"`
}

type ReadBroadcast struct {
	ProfileId   gocql.UUID `cql:"profile_id"`
	BroadcastId gocql.UUID `cql:"broadcast_id"`
}

type UnifiedNotifications struct {
	ID      gocql.UUID `json:"id"`
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Type    string     `json:"type"`
	HasRead bool       `json:"has_read"`
}
type DeliveryStatus struct {
	NotificationId gocql.UUID `cql:"notification_id"`
	ProfileId      gocql.UUID `cql:"profile_id"`
	Channel        string     `cql:"channel"`
	Status         string     `cql:"status"`
	LastAttempt    time.Time  `cql:"last_attempt"`
	Attempts       int        `cql:"attempts"`
}

type BroadcastDeliveryStatus struct {
	BroadcastId gocql.UUID `cql:"broadcast_id"`
	Channel     string     `cql:"channel"`
	Status      string     `cql:"status"`
	Attempts    int        `cql:"attempts"`
	LastAttempt time.Time  `cql:"last_attempt"`
}
