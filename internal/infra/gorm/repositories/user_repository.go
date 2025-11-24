package repositories

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) repositories.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *gormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *gormUserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *gormUserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) FindAllMembers(pagination *models.Pagination) ([]models.User, error) {
	var users []models.User

	query := r.db.Where("role = ?", models.RoleMember).Order(pagination.Sort)

	if pagination.Limit > 0 {
		query = query.Limit(pagination.Limit).Offset(pagination.Offset)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}