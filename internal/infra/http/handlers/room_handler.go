package handlers

import (
	"backend/internal/app/services"
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	roomService services.RoomService
}

func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

// GetAllRooms: Mengambil semua kamar (Public)
func (h *RoomHandler) GetAllRooms(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	pagination := &models.Pagination{
		Page:   page,
		Limit:  limit,
		Sort:   c.Query("sort", "created_at desc"),
		Offset: (page - 1) * limit,
	}

	rooms, err := h.roomService.GetAllRooms(pagination)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mengambil data kamar")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Berhasil mengambil data kamar", fiber.Map{
		"rooms": rooms,
		"page":  page,
		"limit": limit,
	})
}

// GetRoomByID: Mengambil detail kamar (Public)
func (h *RoomHandler) GetRoomByID(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID kamar tidak valid")
	}

	room, err := h.roomService.GetRoomByID(uint(roomID))
	if err != nil {
		if errors.Is(err, errors.New("kamar tidak ditemukan")) {
			return utils.RespondError(c, fiber.StatusNotFound, "Kamar tidak ditemukan")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mengambil data kamar")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Berhasil mengambil data kamar", room)
}

type GetAvailableRoomsInput struct {
	CheckInDate  string `json:"check_in_date" validate:"required"`
	CheckOutDate string `json:"check_out_date" validate:"required"`
}

// GetAvailableRooms: Mengambil kamar yang tersedia (Public)
func (h *RoomHandler) GetAvailableRooms(c *fiber.Ctx) error {
	var input GetAvailableRoomsInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	pagination := &models.Pagination{
		Page:   page,
		Limit:  limit,
		Sort:   c.Query("sort", "created_at desc"),
		Offset: (page - 1) * limit,
	}

	rooms, err := h.roomService.GetAvailableRooms(input.CheckInDate, input.CheckOutDate, pagination)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mengambil kamar tersedia")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Berhasil mengambil kamar tersedia", fiber.Map{
		"rooms": rooms,
		"page":  page,
		"limit": limit,
	})
}

type CreateRoomInput struct {
	RoomNumber   string  `json:"room_number" validate:"required"`
	Type         string  `json:"type" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
	Description  string  `json:"description"`
	MaxOccupancy int     `json:"max_occupancy" validate:"required"`
}

// CreateRoom: Membuat kamar baru (Admin Only)
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var input CreateRoomInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	room := &models.Room{
		RoomNumber:   input.RoomNumber,
		Type:         input.Type,
		Price:        input.Price,
		Description:  input.Description,
		MaxOccupancy: input.MaxOccupancy,
		Status:       "available",
	}

	createdRoom, err := h.roomService.CreateRoom(room)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.RespondSuccess(c, fiber.StatusCreated, "Kamar berhasil dibuat", createdRoom)
}

type UpdateRoomInput struct {
	RoomNumber   string  `json:"room_number"`
	Type         string  `json:"type"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	MaxOccupancy int     `json:"max_occupancy"`
}

// UpdateRoom: Mengubah data kamar (Admin Only)
func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID kamar tidak valid")
	}

	var input UpdateRoomInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	// Ambil room yang ada terlebih dahulu
	existingRoom, err := h.roomService.GetRoomByID(uint(roomID))
	if err != nil {
		return utils.RespondError(c, fiber.StatusNotFound, "Kamar tidak ditemukan")
	}

	// Update field yang diberikan
	if input.RoomNumber != "" {
		existingRoom.RoomNumber = input.RoomNumber
	}
	if input.Type != "" {
		existingRoom.Type = input.Type
	}
	if input.Price > 0 {
		existingRoom.Price = input.Price
	}
	if input.Description != "" {
		existingRoom.Description = input.Description
	}
	if input.Status != "" {
		existingRoom.Status = input.Status
	}
	if input.MaxOccupancy > 0 {
		existingRoom.MaxOccupancy = input.MaxOccupancy
	}

	updatedRoom, err := h.roomService.UpdateRoom(existingRoom)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal mengubah kamar")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Kamar berhasil diubah", updatedRoom)
}

// DeleteRoom: Menghapus kamar (Admin Only)
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID kamar tidak valid")
	}

	if err := h.roomService.DeleteRoom(uint(roomID)); err != nil {
		if errors.Is(err, errors.New("kamar tidak ditemukan")) {
			return utils.RespondError(c, fiber.StatusNotFound, "Kamar tidak ditemukan")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal menghapus kamar")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Kamar berhasil dihapus", nil)
}

type AddRoomImageInput struct {
	ImageURL  string `json:"image_url" validate:"required"`
	IsPrimary bool   `json:"is_primary"`
}

// AddRoomImage: Menambah gambar kamar (Admin Only)
func (h *RoomHandler) AddRoomImage(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID kamar tidak valid")
	}

	var input AddRoomImageInput
	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	roomImage := &models.RoomImage{
		RoomID:    uint(roomID),
		ImageURL:  input.ImageURL,
		IsPrimary: input.IsPrimary,
	}

	createdImage, err := h.roomService.AddRoomImage(roomImage)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.RespondSuccess(c, fiber.StatusCreated, "Gambar kamar berhasil ditambah", createdImage)
}

// DeleteRoomImage: Menghapus gambar kamar (Admin Only)
func (h *RoomHandler) DeleteRoomImage(c *fiber.Ctx) error {
	imageID, err := strconv.ParseUint(c.Params("imageId"), 10, 32)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "ID gambar tidak valid")
	}

	if err := h.roomService.DeleteRoomImage(uint(imageID)); err != nil {
		if errors.Is(err, errors.New("gambar tidak ditemukan")) {
			return utils.RespondError(c, fiber.StatusNotFound, "Gambar tidak ditemukan")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Gagal menghapus gambar")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Gambar kamar berhasil dihapus", nil)
}
