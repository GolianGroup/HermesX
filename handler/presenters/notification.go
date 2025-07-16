package presenters

import "github.com/gofiber/fiber/v2"

func SentNotificationResponse(Title string, Message string) *fiber.Map {
	return SuccessWithDataResponse(fiber.Map{
		"title":   Title,
		"message": Message,
	})
}
