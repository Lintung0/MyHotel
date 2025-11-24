package services

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"errors"

	"gorm.io/gorm"
)

type roomServiceImpl struct {
	roomRepo      repositories.RoomRepository
	roomImageRepo repositories.RoomImageRepository
}

func NewRoomService(rRepo repositories.RoomRepository, riRepo repositories.RoomImageRepository) RoomService {
	return &roomServiceImpl{roomRepo: rRepo, roomImageRepo: riRepo}
}

// GetAllRooms: Mengambil semua kamar dengan pagination
func (s *roomServiceImpl) GetAllRooms(pagination *models.Pagination) ([]models.Room, error) {
	return s.roomRepo.FindAll(pagination)
}

// GetRoomByID: Mengambil detail kamar berdasarkan ID
func (s *roomServiceImpl) GetRoomByID(roomID uint) (*models.Room, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kamar tidak ditemukan")
		}
		return nil, err
	}
	return room, nil
}

// GetAvailableRooms: Mengambil kamar yang tersedia pada periode tertentu
func (s *roomServiceImpl) GetAvailableRooms(checkInDate, checkOutDate string, pagination *models.Pagination) ([]models.Room, error) {
	return s.roomRepo.FindAvailable(checkInDate, checkOutDate, pagination)
}

// CreateRoom: Membuat kamar baru (Admin Only)
func (s *roomServiceImpl) CreateRoom(room *models.Room) (*models.Room, error) {
	// Validasi input
	if room.RoomNumber == "" || room.Type == "" || room.Price <= 0 {
		return nil, errors.New("data kamar tidak lengkap atau tidak valid")
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}
	return room, nil
}

// UpdateRoom: Mengubah data kamar (Admin Only)
func (s *roomServiceImpl) UpdateRoom(room *models.Room) (*models.Room, error) {
	// Verifikasi kamar ada
	_, err := s.roomRepo.FindByID(room.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kamar tidak ditemukan")
		}
		return nil, err
	}

	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}
	return room, nil
}

// DeleteRoom: Menghapus kamar (Admin Only)
func (s *roomServiceImpl) DeleteRoom(roomID uint) error {
	// Verifikasi kamar ada
	_, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kamar tidak ditemukan")
		}
		return err
	}

	// Hapus semua gambar kamar terlebih dahulu
	if err := s.roomImageRepo.DeleteByRoomID(roomID); err != nil {
		return err
	}

	return s.roomRepo.Delete(roomID)
}

// AddRoomImage: Menambah gambar kamar (Admin Only)
func (s *roomServiceImpl) AddRoomImage(image *models.RoomImage) (*models.RoomImage, error) {
	// Verifikasi kamar ada
	_, err := s.roomRepo.FindByID(image.RoomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kamar tidak ditemukan")
		}
		return nil, err
	}

	if err := s.roomImageRepo.Create(image); err != nil {
		return nil, err
	}
	return image, nil
}

// DeleteRoomImage: Menghapus satu gambar kamar (Admin Only)
func (s *roomServiceImpl) DeleteRoomImage(imageID uint) error {
	_, err := s.roomImageRepo.FindByID(imageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gambar tidak ditemukan")
		}
		return err
	}
	return s.roomImageRepo.Delete(imageID)
}

// DeleteRoomImages: Menghapus semua gambar kamar tertentu (Admin Only)
func (s *roomServiceImpl) DeleteRoomImages(roomID uint) error {
	return s.roomImageRepo.DeleteByRoomID(roomID)
}
