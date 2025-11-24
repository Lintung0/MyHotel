package repositories

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type gormRoomRepository struct {
	db *gorm.DB
}

func NewGormRoomRepository(db *gorm.DB) repositories.RoomRepository {
	return &gormRoomRepository{db: db}
}

func (r *gormRoomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *gormRoomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *gormRoomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *gormRoomRepository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	// Preload Images untuk Fitur Galeri Foto
	if err := r.db.Preload("Images").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *gormRoomRepository) FindAll(pagination *models.Pagination) ([]models.Room, error) {
	var rooms []models.Room
	query := r.db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort)
	
	if err := query.Preload("Images").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *gormRoomRepository) FindAvailable(checkInDate, checkOutDate string, pagination *models.Pagination) ([]models.Room, error) {
	var availableRooms []models.Room

	// Subquery untuk mencari Room ID yang sudah dibooking pada periode tertentu
	subQuery := r.db.Model(&models.Booking{}).
		Select("room_id").
		Where("check_out_date > ? AND check_in_date < ?", checkInDate, checkOutDate).
		Where("booking_status IN (?)", []string{models.StatusConfirmed, models.StatusPaid}) 

	// Query utama: Kamar yang ID-nya TIDAK ADA di hasil sub-query, dan statusnya 'available'
	query := r.db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort)
	
	if err := query.Preload("Images").
		Where("id NOT IN (?)", subQuery).
		Where("status = ?", "available").
		Find(&availableRooms).Error; err != nil {
		return nil, err
	}
	return availableRooms, nil
}