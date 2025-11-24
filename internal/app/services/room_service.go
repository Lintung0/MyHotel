package services

import (
	"backend/internal/domain/models"
)

// RoomService mendefinisikan kontrak untuk semua operasi kamar
type RoomService interface {
	// Untuk Member & Admin
	GetAllRooms(pagination *models.Pagination) ([]models.Room, error)
	GetRoomByID(roomID uint) (*models.Room, error)
	GetAvailableRooms(checkInDate, checkOutDate string, pagination *models.Pagination) ([]models.Room, error)

	// Untuk Admin
	CreateRoom(room *models.Room) (*models.Room, error)
	UpdateRoom(room *models.Room) (*models.Room, error)
	DeleteRoom(roomID uint) error

	// Untuk Galeri Foto
	AddRoomImage(image *models.RoomImage) (*models.RoomImage, error)
	DeleteRoomImage(imageID uint) error
	DeleteRoomImages(roomID uint) error
}
