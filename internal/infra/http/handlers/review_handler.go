package handlers

import (
	"backend/internal/app/services"
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	reviewService services.ReviewService
}

func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

type CreateReviewInput struct {
	BookingID uint   `json:"booking_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}

// CreateReview: Membuat review untuk booking (Member Only)
func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input CreateReviewInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	review := &models.Review{
		BookingID: input.BookingID,
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	createdReview, err := h.reviewService.CreateReview(review)
	if err != nil {
		if errors.Is(err, errors.New("anda sudah memberikan ulasan untuk pemesanan ini")) {
			return utils.RespondError(c, fiber.StatusConflict, "Anda sudah memberikan ulasan untuk pemesanan ini")
		}
		return utils.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.RespondSuccess(c, fiber.StatusCreated, "Ulasan berhasil dibuat", createdReview)
}

// GetRoomReviews: Mengambil semua review kamar (Public)
func (h *ReviewHandler) GetRoomReviews(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("roomId"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID kamar tidak valid")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	pagination := &models.Pagination{
		Page:   page,
		Limit:  limit,
		Sort:   c.Query("sort", "created_at desc"),
		Offset: (page - 1) * limit,
	}

	reviews, err := h.reviewService.GetRoomReviews(uint(roomID), pagination)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mengambil data ulasan")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Berhasil mengambil data ulasan", fiber.Map{
		"reviews": reviews,
		"page":    page,
		"limit":   limit,
	})
}

// GetReviewByID: Mengambil detail review (Public)
func (h *ReviewHandler) GetReviewByID(c *fiber.Ctx) error {
	reviewID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID ulasan tidak valid")
	}

	review, err := h.reviewService.GetReviewByID(uint(reviewID))
	if err != nil {
		if errors.Is(err, errors.New("ulasan tidak ditemukan")) {
			return utils.RespondError(c, fiber.StatusNotFound, "Ulasan tidak ditemukan")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mengambil data ulasan")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Berhasil mengambil data ulasan", review)
}

// DeleteReview: Menghapus review (Admin Only)
func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	reviewID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID ulasan tidak valid")
	}

	if err := h.reviewService.DeleteReview(uint(reviewID)); err != nil {
		if errors.Is(err, errors.New("ulasan tidak ditemukan")) {
			return utils.RespondError(c, fiber.StatusNotFound, "Ulasan tidak ditemukan")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal menghapus ulasan")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Ulasan berhasil dihapus", nil)
}
