package handlers

import (
	"backend/internal/app/services"
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	token, user, err := h.authService.Login(input.Username, input.Password)

	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) || errors.Is(err, models.ErrInvalidCredentials) {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Username atau Password Salah")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal login")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Login Berhasil", fiber.Map{
		"token": token,
		"user":  user,
	})
}

type RegisterInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	newUser := &models.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
		FullName: input.FullName,
		Role:     models.RoleMember, // Default: Member
	}

	_, err := h.authService.Register(newUser)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return utils.RespondError(c, fiber.StatusConflict, "Username atau email sudah digunakan")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mendaftarkan user")
	}
	return utils.RespondSuccess(c, fiber.StatusOK, "Pendaftaran berhasil", nil)
}
