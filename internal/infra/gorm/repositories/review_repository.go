package repositories

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type gormReviewRepository struct {
	db *gorm.DB
}

func NewGormReviewRepository(db *gorm.DB) repositories.ReviewRepository {
	return &gormReviewRepository{db: db}
}

func (r *gormReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *gormReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *gormReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *gormReviewRepository) FindByID(id uint) (*models.Review, error) {
	var review models.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *gormReviewRepository) FindByBookingID(bookingID uint) (*models.Review, error) {
	var review models.Review
	if err := r.db.Where("booking_id = ?", bookingID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *gormReviewRepository) FindByRoomID(roomID uint, pagination *models.Pagination) ([]models.Review, error) {
	var reviews []models.Review
    
	query := r.db.
        // Joins Booking untuk mendapatkan RoomID
        Joins("JOIN bookings ON bookings.id = reviews.booking_id").
        Where("bookings.room_id = ?", roomID).
        Order(pagination.Sort)

	if pagination.Limit > 0 {
		query = query.Limit(pagination.Limit).Offset(pagination.Offset)
	}
    
    // Preload User untuk menampilkan nama reviewer
	if err := query.Preload("User").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}