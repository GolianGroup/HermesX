package presenters

import (
	"hermesx/internal/repositories/models"

	"github.com/gofiber/fiber/v2"
)

func SentNotificationResponse(Title string, Message string) *fiber.Map {
	return SuccessWithDataResponse(fiber.Map{
		"title":   Title,
		"message": Message,
	})
}

type Notification struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
	HasRead bool   `json:"has_read"`
}

func UnifiedNotificationsResponse(notifications []models.UnifiedNotifications) *fiber.Map {
	cards := make([]Notification, len(notifications))
	for i, model := range notifications {
		cards[i] = Notification{
			ID:      model.ID.String(),
			Title:   model.Title,
			Message: model.Message,
			Type:    model.Type,
			HasRead: model.HasRead,
		}
	}

	return SuccessWithDataResponse(cards)
}
