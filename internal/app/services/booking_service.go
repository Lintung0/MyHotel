package services

import "backend/internal/domain/models"

// BookingService mendefinisikan kontrak untuk semua operasi pemesanan
type BookingService interface {
	// Untuk Member
	CreateBooking(booking *models.Booking) (*models.Booking, error)
	GetUserBookings(userID uint, pagination *models.Pagination) ([]models.Booking, error)
	CancelBooking(bookingID uint, userID uint) error
	
	// Untuk Admin
	GetAllBookings(pagination *models.Pagination) ([]models.Booking, error)
	UpdatePaymentStatus(bookingID uint, newStatus string) (*models.Booking, error)
	
	// Fitur Review/Ulasan (setelah booking selesai)
	CreateReview(review *models.Review) (*models.Review, error)
}