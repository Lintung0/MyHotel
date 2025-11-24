package middleware

import (
	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	CtxUserIDKey = "user_id"
	CtxRoleKey   = "user_role"
)

// JWTMiddleware: Validasi JWT Token
func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Token tidak ditemukan")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Format token tidak valid (gunakan 'Bearer <token>')")
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Token tidak valid atau sudah kadaluarsa")
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Token tidak dapat diproses")
		}

		c.Locals(CtxUserIDKey, claims.UserID)
		c.Locals(CtxRoleKey, claims.Role)
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware: Cek role user
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals(CtxRoleKey).(string)

		allowed := false
		for _, r := range allowedRoles {
			if role == r {
				allowed = true
				break
			}
		}

		if !allowed {
			return utils.RespondError(c, fiber.StatusForbidden, "Anda tidak memiliki akses ke resource ini")
		}

		return c.Next()
	}
}

// JWTProtected: Legacy function untuk backward compatibility
func JWTProtected(cfg *config.Config) fiber.Handler {
	return JWTMiddleware(cfg)
}
