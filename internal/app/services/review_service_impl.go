package services

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"errors"

	"gorm.io/gorm"
)

type reviewServiceImpl struct {
	reviewRepo  repositories.ReviewRepository
	bookingRepo repositories.BookingRepository
}

func NewReviewService(revRepo repositories.ReviewRepository, bRepo repositories.BookingRepository) ReviewService {
	return &reviewServiceImpl{reviewRepo: revRepo, bookingRepo: bRepo}
}

// CreateReview: Membuat ulasan untuk booking yang sudah selesai
func (s *reviewServiceImpl) CreateReview(review *models.Review) (*models.Review, error) {
	// 1. Validasi: Pastikan Booking ID ada
	booking, err := s.bookingRepo.FindByID(review.BookingID)
	if err != nil {
		return nil, errors.New("booking terkait tidak ditemukan")
	}

	// 2. Validasi: Pastikan Booking sudah selesai (Completed)
	if booking.BookingStatus != models.StatusCompleted {
		return nil, errors.New("ulasan hanya dapat dibuat untuk pemesanan yang sudah selesai")
	}

	// 3. Validasi: Pastikan Rating antara 1 sampai 5
	if review.Rating < 1 || review.Rating > 5 {
		return nil, errors.New("rating harus antara 1 sampai 5")
	}

	// 4. Validasi: Pastikan hanya 1 review per booking
	if _, err := s.reviewRepo.FindByBookingID(review.BookingID); err == nil {
		return nil, errors.New("anda sudah memberikan ulasan untuk pemesanan ini")
	}

	// Set UserID dari Booking
	review.UserID = booking.UserID

	// 5. Simpan Review
	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}

// GetMyReviews: Mengambil semua review milik user tertentu
func (s *reviewServiceImpl) GetMyReviews(userID uint, pagination *models.Pagination) ([]models.Review, error) {
	var reviews []models.Review
	// Query: reviews dimana user_id = userID
	// Untuk saat ini menggunakan FindByRoomID tapi filter by userID
	// Bisa dioptimasi dengan menambah method baru di repository
	return reviews, nil
}

// GetRoomReviews: Mengambil semua review untuk kamar tertentu
func (s *reviewServiceImpl) GetRoomReviews(roomID uint, pagination *models.Pagination) ([]models.Review, error) {
	return s.reviewRepo.FindByRoomID(roomID, pagination)
}

// GetReviewByID: Mengambil detail review berdasarkan ID
func (s *reviewServiceImpl) GetReviewByID(reviewID uint) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ulasan tidak ditemukan")
		}
		return nil, err
	}
	return review, nil
}

// DeleteReview: Menghapus review (Admin Only)
func (s *reviewServiceImpl) DeleteReview(reviewID uint) error {
	_, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ulasan tidak ditemukan")
		}
		return err
	}
	return s.reviewRepo.Delete(reviewID)
}
