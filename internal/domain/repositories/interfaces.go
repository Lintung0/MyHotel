package repositories

import (
	"backend/internal/domain/models"
)

type RoomRepository interface {
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(id uint) error
	FindByID(id uint) (*models.Room, error)

	// Show & Search
	FindAll(pagination *models.Pagination) ([]models.Room, error)
	// Fungsi untuk Filter Ketersediaan Real-time
	FindAvailable(checkInDate, checkOutDate string, pagination *models.Pagination) ([]models.Room, error)
}

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	// Tambahan untuk Admin
	FindAllMembers(pagination *models.Pagination) ([]models.User, error)
}

type BookingRepository interface {
	Create(booking *models.Booking) error
	Update(booking *models.Booking) error
	Delete(id uint) error
	FindByID(id uint) (*models.Booking, error)

	// Fungsi Member dan Admin
	FindByUserID(userID uint, pagination *models.Pagination) ([]models.Booking, error)
	FindAll(pagination *models.Pagination) ([]models.Booking, error) // Untuk Admin melihat semua

	// Fungsi Logika Bisnis
	UpdateStatus(id uint, newStatus string) error                             // Mengubah booking/payment status oleh Admin
	CheckOverlap(roomID uint, checkInDate, checkOutDate string) (bool, error) // Pencegahan Double Booking
}

type RoomImageRepository interface {
	Create(image *models.RoomImage) error
	Update(image *models.RoomImage) error
	Delete(id uint) error
	FindByID(id uint) (*models.RoomImage, error)
	FindByRoomID(roomID uint) ([]models.RoomImage, error)
	// Tambahan
	DeleteByRoomID(roomID uint) error
}

type ReviewRepository interface {
	Create(review *models.Review) error
	Update(review *models.Review) error
	Delete(id uint) error
	FindByID(id uint) (*models.Review, error)
	FindByBookingID(bookingID uint) (*models.Review, error)
	// Tambahan untuk tampilan kamar
	FindByRoomID(roomID uint, pagination *models.Pagination) ([]models.Review, error)
}
