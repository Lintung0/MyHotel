package mysql

import (
	"backend/internal/domain/models"

	"gorm.io/gorm"
)

type RoomRepositoryImpl struct {
	DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepositoryImpl {
	return &RoomRepositoryImpl{DB: db}
}

// create room (insert)
func (r *RoomRepositoryImpl) Create(room *models.Room) error {
	return r.DB.Create(room).Error
}

// update room (update)
func (r *RoomRepositoryImpl) Update(room *models.Room) error {
	return r.DB.Model(room).Updates(room).Error
}

// Delete room (delete)
func (r *RoomRepositoryImpl) Delete(id uint) error {
	return r.DB.Delete(&models.Room{}, id).Error
}

// Find room by ID (select by id)
func (r *RoomRepositoryImpl) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := r.DB.Preload("Images").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// FindAll mengimplementasikan penampilan semua kamar dengan Pagination
func (r *RoomRepositoryImpl) FindAll(pagination *models.Pagination) ([]models.Room, error) {
	var rooms []models.Room

	query := r.DB.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort)

	if err := query.Preload("Images").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomRepositoryImpl) FindAvailable(checkInDate, checkOutDate string, pagination *models.Pagination) ([]models.Room, error) {
	var availableRooms []models.Room

	subQuery := r.DB.Model(&models.Booking{}).
		Select("room_id").
		Where("check_out_date > ? AND check_in_date > ?", checkInDate, checkOutDate)

	query := r.DB.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort)

	if err := query.Preload("Images").
		Where("id NOT IN (?)", subQuery).
		Where("status = ?", "available").
		Find(&availableRooms).Error; err != nil {
		return nil, err
	}
	return availableRooms, nil
}
