package middleware

import (
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthorizeRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals(CtxRoleKey).(string)
		if !ok {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Akses ditolak: user belum terauthentikasi")
		}

		if userRole != requiredRole {
			return utils.RespondError(c, fiber.StatusForbidden, "Akses ditolak: anda tidak memiliki hak akses "+requiredRole)
		}

		return c.Next()
	}
}
