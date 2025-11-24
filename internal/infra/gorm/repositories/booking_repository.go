package repositories

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"errors"

	"gorm.io/gorm"
)

type gormBookingRepository struct {
	db *gorm.DB
}

func NewGormBookingRepository(db *gorm.DB) repositories.BookingRepository {
	return &gormBookingRepository{db: db}
}

func (r *gormBookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *gormBookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *gormBookingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Booking{}, id).Error
}

func (r *gormBookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.Preload("Room").Preload("User").First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *gormBookingRepository) FindByUserID(userID uint, pagination *models.Pagination) ([]models.Booking, error) {
	var bookings []models.Booking
	query := r.db.Where("user_id = ?", userID).Order(pagination.Sort)

	if pagination.Limit > 0 {
		query = query.Limit(pagination.Limit).Offset(pagination.Offset)
	}

	if err := query.Preload("Room").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *gormBookingRepository) FindAll(pagination *models.Pagination) ([]models.Booking, error) {
	var bookings []models.Booking
	query := r.db.Order(pagination.Sort)

	if pagination.Limit > 0 {
		query = query.Limit(pagination.Limit).Offset(pagination.Offset)
	}

	if err := query.Preload("Room").Preload("User").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *gormBookingRepository) UpdateStatus(id uint, newStatus string) error {
	// Mengupdate BookingStatus (logika PaymentStatus akan dihandle di Service Layer)
	result := r.db.Model(&models.Booking{}).Where("id = ?", id).Update("booking_status", newStatus)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("booking tidak ditemukan atau status tidak berubah")
	}
	return nil
}

func (r *gormBookingRepository) CheckOverlap(roomID uint, checkInDate, checkOutDate string) (bool, error) {
	var count int64

	err := r.db.Model(&models.Booking{}).
		Where("room_id = ?", roomID).
		Where("booking_status IN (?)", []string{models.StatusConfirmed, models.StatusPaid}).
		Where("check_out_date > ? AND check_in_date < ?", checkInDate, checkOutDate).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
