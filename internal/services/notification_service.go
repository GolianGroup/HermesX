package services

import (
	"context"
	"encoding/json"
	"fmt"
	"hermesx/handler/dtos"
	"hermesx/internal/helper/nats"
	"hermesx/internal/repositories"
	"hermesx/internal/repositories/models"
	"time"

	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

type NotificationService interface {
	SendNotification(ctx context.Context, notification *dtos.Notification) error
	BroadcastNotification(ctx context.Context, broadcast *dtos.Broadcast) error
	ReadNotification(ctx context.Context, notificationId gocql.UUID, profileId gocql.UUID) error
	ReadBroadcastNotification(ctx context.Context, broadcastId gocql.UUID, profileId gocql.UUID) error
	GetUserNotifications(ctx context.Context, profileId gocql.UUID) ([]models.UnifiedNotifications, error)
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
	eventRepo        repositories.EventRepository
	profileRepo      repositories.ProfileRepository
	logger           *zap.Logger
	nats             nats.NatsConnection
}

func NewNotificationService(
	notificationRepo repositories.NotificationRepository,
	eventRepo repositories.EventRepository,
	profileRepo repositories.ProfileRepository,
	logger *zap.Logger,
	nats nats.NatsConnection,
) NotificationService {
	notificationServiceLogger := logger.With(zap.String("service", "notification"))
	return &notificationService{
		notificationRepo: notificationRepo, eventRepo: eventRepo, profileRepo: profileRepo, logger: notificationServiceLogger, nats: nats,
	}
}

func (n *notificationService) SendNotification(ctx context.Context, notification *dtos.Notification) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id := gocql.TimeUUID()

	model := models.Notification{
		ID:        id,
		ProfileId: notification.ProfileId,
		Title:     notification.Title,
		Message:   notification.Message,
	}

	event, err := n.eventRepo.GetByName(ctx, notification.Event)
	if err != nil {
		n.logger.Debug("Error in SendNotification", zap.Error(err))
		return err
	}

	profile, err := n.profileRepo.GetProfile(ctx, notification.ProfileId.String(), notification.ProfileId.String())
	if err != nil {
		n.logger.Debug("failed to get profile", zap.Error(err), zap.String("profile_id", notification.ProfileId.String()))
		return err
	}

	// send Critical Message
	if event.IsCritical {
		subject := fmt.Sprintf("individual.%v", event.PreferredChannel)

		n.logger.Debug("Sendning Critical Notification to", zap.String("profile_id", notification.ProfileId.String()), zap.String("channel", subject))
		deliveryModel := models.DeliveryStatus{
			NotificationId: id,
			ProfileId:      model.ProfileId,
			Event:          event.Name,
			Status:         "initiated",
			Attempts:       1,
			LastAttempt:    time.Now().UTC(),
		}
		err := n.notificationRepo.SetDeliveryStatus(ctx, deliveryModel)
		if err != nil {
			n.logger.Error("Failed to save delivery status", zap.Error(err))
			return err
		}
		natsMessage := nats.CriticalMessage{
			Type:     event.PreferredChannel,
			Messsage: notification.Message,
		}
		switch event.PreferredChannel {
		case "sms":
			natsMessage.Destination = profile.PhoneNumber
		case "email":
			natsMessage.Destination = profile.Email
			natsMessage.Template = event.Template
		}

		if err := n.publishCriticalNotification(ctx, subject, natsMessage); err != nil {
			deliveryModel.Status = "failed"
			err = n.notificationRepo.SetDeliveryStatus(ctx, deliveryModel)
			if err != nil {
				n.logger.Error("Failed to update delivery status", zap.Error(err))
				return err
			}
			return nil
		}

		deliveryModel.Status = "delivered"
		if err := n.notificationRepo.SetDeliveryStatus(ctx, deliveryModel); err != nil {
			n.logger.Error("Failed to update delivery status", zap.Error(err))
			return err
		}
		if !event.IsProtected {
			err = n.notificationRepo.CreateNotification(ctx, model)
			if err != nil {
				n.logger.Debug("failed to create critical notification", zap.Error(err))
				return err
			}
		}
		n.logger.Info("Critical notification delivered", zap.String("notification_id", model.ID.String()))
		return nil

	}

	subject := fmt.Sprintf("individual.%v", event.PreferredChannel)
	natsMessage := nats.NormalMessage{
		Type:     event.PreferredChannel,
		Messsage: notification.Message,
	}

	switch event.PreferredChannel {
	case "sms":
		if profile.Preferences.Notifications.SMS {
			natsMessage.Destination = profile.PhoneNumber
		} else {
			n.logger.Debug("Notification ignored during user prefrences for sms",
				zap.String("profile_id", notification.ProfileId.String()),
				zap.String("event", event.Name),
				zap.String("subject", subject),
			)
			return nil
		}
	case "email":
		if profile.Preferences.Notifications.Email {
			natsMessage.Destination = profile.Email
			natsMessage.Type = event.Template
		} else {
			n.logger.Debug("Notification ignored during user prefrences for email",
				zap.String("profile_id", notification.ProfileId.String()),
				zap.String("event", event.Name),
				zap.String("subject", subject),
			)
			return nil
		}
	}

	n.logger.Debug("Sendning Notification to",
		zap.String("profile_id", notification.ProfileId.String()),
		zap.String("event", event.Name),
		zap.String("subject", subject),
	)

	err = n.notificationRepo.CreateNotification(ctx, model)
	if err != nil {
		n.logger.Debug("failed to create notification", zap.Error(err))
		return err
	}

	err = n.publishNotification(ctx, subject, natsMessage)
	if err != nil {
		n.logger.Error("Failed to send notification",
			zap.String("event", event.Name),
			zap.String("subject", subject),
			zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationService) ReadNotification(ctx context.Context, notificationId gocql.UUID, profileId gocql.UUID) error {
	err := n.notificationRepo.ReadNotification(ctx, notificationId, profileId)
	if err != nil {
		n.logger.Error("Failed to read notification", zap.String("profile_id", profileId.String()), zap.String("notification_id", notificationId.String()))
		return err
	}
	return nil
}

// This function only works for push notifications TODO: must develop a workerpool to schedule broadcast emails and sms
func (n *notificationService) BroadcastNotification(ctx context.Context, broadcast *dtos.Broadcast) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id := gocql.TimeUUID()

	model := models.Broadcast{
		ID:      id,
		Title:   broadcast.Title,
		Message: broadcast.Message,
	}
	event, err := n.eventRepo.GetByName(ctx, broadcast.Event)
	if err != nil {
		n.logger.Debug("Error in BroadcastNotification", zap.Error(err))
		return err
	}

	if event.IsCritical {
		subject := fmt.Sprintf("broadcast.%v", event.PreferredChannel)

		n.logger.Debug("Broadcast Critical Notification")
		deliverModel := models.BroadcastDeliveryStatus{
			BroadcastId: id,
			Event:       event.Name,
			Status:      "initiated",
			Attempts:    1,
			LastAttempt: time.Now().UTC(),
		}
		err := n.notificationRepo.SetBroadcastDeliveryStatus(ctx, deliverModel)
		if err != nil {
			n.logger.Error("Failed to save broadcast delivery status", zap.Error(err))
			return err
		}

		natsMessage := nats.Broadcast{
			Title:    broadcast.Title,
			Messsage: broadcast.Message,
		}
		if err := n.publishBroadcast(ctx, subject, natsMessage); err != nil {
			deliverModel.Status = "failed"
			err = n.notificationRepo.SetBroadcastDeliveryStatus(ctx, deliverModel)
			if err != nil {
				n.logger.Error("Failed to update delivery status", zap.Error(err))
				return err
			}
			return nil
		}
		deliverModel.Status = "delivered"
		if err := n.notificationRepo.SetBroadcastDeliveryStatus(ctx, deliverModel); err != nil {
			n.logger.Error("Failed to update broadcast delivery status", zap.Error(err))
			return err
		}
		n.logger.Info("Critical broadcast delivered", zap.String("broadcast_id", model.ID.String()))
		return nil
	}

	subject := fmt.Sprintf("broadcast.%v", event.PreferredChannel)
	n.logger.Debug("Sendning Broadcast notification in ", zap.String("event", subject))
	err = n.notificationRepo.CreateBroadcastNotification(ctx, model)
	if err != nil {
		n.logger.Debug("failed to create broadcast")
		return err
	}

	natsMessage := nats.Broadcast{
		Title:    broadcast.Title,
		Messsage: broadcast.Message,
	}
	err = n.publishBroadcast(ctx, subject, natsMessage)
	if err != nil {
		n.logger.Error("Failed to send broadcast on channel", zap.String("channel", subject), zap.Error(err))
		return err
	}

	return nil
}

func (n *notificationService) ReadBroadcastNotification(ctx context.Context, broadcastId gocql.UUID, profileId gocql.UUID) error {
	readModel := models.ReadBroadcast{
		ProfileId:   profileId,
		BroadcastId: broadcastId,
	}

	err := n.notificationRepo.CreateReadBroadcastNotification(ctx, readModel)
	if err != nil {
		n.logger.Error("Failed to read notification", zap.String("profile_id", readModel.ProfileId.String()), zap.String("notification_id", readModel.BroadcastId.String()))
		return err
	}
	return nil
}

func (n *notificationService) GetUserNotifications(ctx context.Context, profileId gocql.UUID) ([]models.UnifiedNotifications, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var result []models.UnifiedNotifications

	userNotifs, err := n.notificationRepo.FetchUserNotifications(ctx, profileId)
	if err != nil {
		n.logger.Error("Failed to fetch user notifications", zap.Error(err))
		return nil, err
	}
	for _, notifs := range userNotifs {
		result = append(result, models.UnifiedNotifications{
			ID:      notifs.ID,
			Title:   notifs.Title,
			Message: notifs.Message,
			Type:    "individual",
			HasRead: notifs.HasRead,
		})
	}

	broadcastNotifs, err := n.notificationRepo.FetchBroadcasts(ctx)
	if err != nil {
		n.logger.Error("Failed to fetch broadcast notifications", zap.Error(err))
		return nil, err
	}

	broadcastReads, err := n.notificationRepo.FetchUserReadBroadcast(ctx, profileId)
	if err != nil {
		n.logger.Error("Failed to fetch broadcast notifications", zap.Error(err))
		return nil, err
	}

	readSet := make(map[gocql.UUID]struct{}, len(broadcastReads))
	for _, id := range broadcastReads {
		readSet[id] = struct{}{}
	}

	for _, notif := range broadcastNotifs {
		_, hasRead := readSet[notif.ID]

		result = append(result, models.UnifiedNotifications{
			ID:      notif.ID,
			Title:   notif.Title,
			Message: notif.Message,
			Type:    "broadcast",
			HasRead: hasRead,
		})
	}

	return result, nil

}

func (n *notificationService) publishCriticalNotification(_ context.Context, subject string, message nats.CriticalMessage) error {
	// Convert message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		n.logger.Debug("Failed to marshal critical message", zap.Error(err))
		return err
	}

	// Publish the message
	if err := n.nats.Publish(subject, messageBytes); err != nil {
		n.logger.Debug("Failed to publish notification",
			zap.Error(err),
			zap.String("destination", message.Destination),
		)
		return err
	}
	return nil
}

func (n *notificationService) publishNotification(_ context.Context, subject string, message nats.NormalMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		n.logger.Debug("Failed to marshal critical message", zap.Error(err))
		return err
	}

	// Publish the message
	if err := n.nats.Publish(subject, messageBytes); err != nil {
		n.logger.Debug("Failed to publish like/dislike message",
			zap.Error(err),
			zap.String("destination", message.Destination),
		)
		return err
	}
	return nil
}

func (n *notificationService) publishBroadcast(_ context.Context, subject string, message nats.Broadcast) error {
	// Convert message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		n.logger.Debug("Failed to marshal message", zap.Error(err))
		return err
	}

	// Publish the message
	if err := n.nats.Publish(subject, messageBytes); err != nil {
		n.logger.Debug("Failed to publish broadcast notification",
			zap.Error(err),
		)
		return err
	}
	return nil
}

// TODO: Write a function to Broadcast sms and email
