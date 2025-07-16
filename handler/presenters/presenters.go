package presenters

import "github.com/gofiber/fiber/v2"

func SuccessWithDataResponse(data interface{}) *fiber.Map {
	return &fiber.Map{
		"status": "success",
		"data":   data,
	}
}
