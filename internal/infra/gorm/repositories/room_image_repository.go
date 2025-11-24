package repositories

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type gormRoomImageRepository struct {
	db *gorm.DB
}

func NewGormRoomImageRepository(db *gorm.DB) repositories.RoomImageRepository {
	return &gormRoomImageRepository{db: db}
}

func (r *gormRoomImageRepository) Create(image *models.RoomImage) error {
	return r.db.Create(image).Error
}

func (r *gormRoomImageRepository) Update(image *models.RoomImage) error {
	return r.db.Save(image).Error
}

func (r *gormRoomImageRepository) Delete(id uint) error {
	return r.db.Delete(&models.RoomImage{}, id).Error
}

func (r *gormRoomImageRepository) FindByID(id uint) (*models.RoomImage, error) {
	var image models.RoomImage
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *gormRoomImageRepository) FindByRoomID(roomID uint) ([]models.RoomImage, error) {
	var images []models.RoomImage
	if err := r.db.Where("room_id = ?", roomID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *gormRoomImageRepository) DeleteByRoomID(roomID uint) error {
	// Hapus secara permanen semua RoomImage yang terasosiasi dengan RoomID
	return r.db.Unscoped().Where("room_id = ?", roomID).Delete(&models.RoomImage{}).Error
}