package services

import "backend/internal/domain/models"

// ReviewService mendefinisikan kontrak untuk semua operasi review/ulasan
type ReviewService interface {
	// Untuk Member
	CreateReview(review *models.Review) (*models.Review, error)
	GetMyReviews(userID uint, pagination *models.Pagination) ([]models.Review, error)

	// Untuk Public & Admin
	GetRoomReviews(roomID uint, pagination *models.Pagination) ([]models.Review, error)
	GetReviewByID(reviewID uint) (*models.Review, error)

	// Untuk Admin
	DeleteReview(reviewID uint) error
}
