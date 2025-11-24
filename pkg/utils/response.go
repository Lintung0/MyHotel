package utils

import "github.com/gofiber/fiber/v2"

// Struktur standar untuk response API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RespondSuccess mengirim response sukses
func RespondSuccess(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondError mengirim response error
func RespondError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
