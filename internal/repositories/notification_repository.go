package repositories

import (
	"context"
	"fmt"
	"hermesx/internal/database/scylla"
	"hermesx/internal/repositories/models"
	"time"

	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification models.Notification) error
	CreateBroadcastNotification(ctx context.Context, notification models.Broadcast) error
	CreateReadBroadcastNotification(ctx context.Context, read models.ReadBroadcast) error
	ReadNotification(ctx context.Context, id, profileId gocql.UUID) error
	SetDeliveryStatus(ctx context.Context, notification models.DeliveryStatus) error
	SetBroadcastDeliveryStatus(ctx context.Context, broadcast models.BroadcastDeliveryStatus) error
	FetchUserNotifications(ctx context.Context, profileId gocql.UUID) ([]models.GetNotification, error)
	FetchBroadcasts(ctx context.Context) ([]models.Broadcast, error)
	FetchUserReadBroadcast(ctx context.Context, profileId gocql.UUID) ([]gocql.UUID, error)
}

type notificationRepository struct {
	scylla scylla.ScyllaDB
	logger *zap.Logger
}

func NewNotificationRepository(scylla scylla.ScyllaDB, logger *zap.Logger) NotificationRepository {
	return &notificationRepository{scylla: scylla, logger: logger}
}

func (n *notificationRepository) CreateNotification(ctx context.Context, notification models.Notification) error {

	query := `INSERT INTO notification (id, profile_id, has_read, title, message, created_at) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := n.scylla.Session().
		Query(
			query,
			notification.ID,
			notification.ProfileId,
			false,
			notification.Title,
			notification.Message,
			time.Now().UTC(),
		).
		WithContext(ctx).
		Consistency(gocql.One).
		ScanCAS(nil, nil, nil, nil, nil, nil)

	if err != nil {
		n.logger.Debug("Error in SendNotification", zap.Error(err))
		return err
	}
	return nil
}

func (n *notificationRepository) CreateBroadcastNotification(ctx context.Context, notification models.Broadcast) error {
	query := `INSERT INTO broadcast (id, title, message, created_at) VALUES (?, ?, ?, ?) IF NOT EXISTS`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := n.scylla.Session().
		Query(
			query,
			notification.ID,
			notification.Title,
			notification.Message,
			time.Now().UTC(),
		).
		WithContext(ctx).
		Consistency(gocql.One).
		ScanCAS(nil, nil, nil, nil)

	if err != nil {
		n.logger.Debug("Error in BroadcastNotification", zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationRepository) CreateReadBroadcastNotification(ctx context.Context, read models.ReadBroadcast) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	query := `INSERT INTO broadcast_read (profile_id, broadcast_id, read_at) VALUES (?, ?, ?) IF NOT EXISTS`
	_, err := n.scylla.Session().
		Query(
			query,
			read.ProfileId,
			read.BroadcastId,
			time.Now().UTC(),
		).
		WithContext(ctx).
		Consistency(gocql.One).
		ScanCAS(nil, nil, nil)

	if err != nil {
		n.logger.Debug("Error in ReadBroadcastNotification", zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationRepository) ReadNotification(ctx context.Context, id, profileId gocql.UUID) error {
	query := `UPDATE notification SET has_read = 1 WHERE profile_id = ? AND id = ?`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := n.scylla.Session().
		Query(
			query,
			profileId,
			id,
		).
		WithContext(ctx).
		Consistency(gocql.One).
		Exec()

	if err != nil {
		n.logger.Debug("Error in ReadNotification", zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationRepository) SetDeliveryStatus(ctx context.Context, notification models.DeliveryStatus) error {
	query := `UPDATE notification_delivery_status SET status = ?, attempts = ?, event = ?, last_attempt = ? WHERE profile_id = ? AND notification_id = ?`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := n.scylla.Session().Query(
		query,
		notification.Status,
		notification.Attempts,
		notification.Event,
		notification.LastAttempt,
		notification.ProfileId,
		notification.NotificationId,
	).WithContext(ctx).Consistency(gocql.One).Exec()
	if err != nil {
		n.logger.Debug("Error in SetDeliveryStatus", zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationRepository) SetBroadcastDeliveryStatus(ctx context.Context, broadcast models.BroadcastDeliveryStatus) error {
	query := `UPDATE broadcast_delivery_status SET status = ?, attempts = ?, event = ?, last_attempt = ? WHERE broadcast_id = ?`
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := n.scylla.Session().Query(
		query,
		broadcast.Status,
		broadcast.Attempts,
		broadcast.Event,
		broadcast.LastAttempt,
		broadcast.BroadcastId,
	).WithContext(ctx).Consistency(gocql.One).Exec()
	if err != nil {
		n.logger.Debug("Error in SetBroadcastDeliveryStatus", zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationRepository) FetchUserNotifications(ctx context.Context, profileId gocql.UUID) ([]models.GetNotification, error) {
	query := `SELECT id, title, message, has_read FROM notification WHERE profile_id = ?;`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	iter := n.scylla.Session().
		Query(query, profileId).
		WithContext(ctx).
		Iter()

	var notifications []models.GetNotification
	var id gocql.UUID
	var title, message string
	var hasRead bool

	for iter.Scan(&id, &title, &message, &hasRead) {
		notifications = append(notifications, models.GetNotification{
			ID:      id,
			Title:   title,
			Message: message,
			HasRead: hasRead,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, nil

}

func (n *notificationRepository) FetchBroadcasts(ctx context.Context) ([]models.Broadcast, error) {
	query := `SELECT id, title, message FROM broadcast;`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	iter := n.scylla.Session().
		Query(query).
		WithContext(ctx).
		Iter()

	var notifications []models.Broadcast
	var id gocql.UUID
	var title, message string

	for iter.Scan(&id, &title, &message) {
		notifications = append(notifications, models.Broadcast{
			ID:      id,
			Title:   title,
			Message: message,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, nil

}

func (n *notificationRepository) FetchUserReadBroadcast(ctx context.Context, profileId gocql.UUID) ([]gocql.UUID, error) {
	query := `SELECT broadcast_id FROM broadcast_read WHERE profile_id = ?;`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	iter := n.scylla.Session().
		Query(query, profileId).
		WithContext(ctx).
		Iter()

	var ids []gocql.UUID
	var broadcast_id gocql.UUID

	for iter.Scan(&broadcast_id) {
		ids = append(ids, broadcast_id)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	return ids, nil
}
