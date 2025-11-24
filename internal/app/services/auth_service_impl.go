package services

import (
	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implementasi dari interface AuthService
type authServiceImpl struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

// NewAuthService adalah constructor
func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authServiceImpl{userRepo: userRepo, cfg: cfg}
}

// Helper: generateToken membuat JWT Token
func (s *authServiceImpl) generateToken(user *models.User) (string, error) {
	// Waktu kedaluwarsa token
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour)

	// Payload/Claims Token (data user yang disimpan)
	claims := models.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Menandatangani token dengan Secret Key dari .env
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// Register melakukan hashing dan menyimpan user ke DB
func (s *authServiceImpl) Register(user *models.User) (*models.User, error) {
	// 1. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// 2. Set Role Default (Jika tidak diset, GORM akan menggunakan default 'member')
	if user.Role == "" {
		user.Role = models.RoleMember
	}

	// 3. Simpan ke Database melalui Repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Sembunyikan password sebelum dikembalikan
	user.Password = ""
	return user, nil
}

// Login memverifikasi user, password, dan membuat token
func (s *authServiceImpl) Login(username, password string) (string, *models.User, error) {
	// 1. Cari User di DB berdasarkan username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// Gunakan error gorm.ErrRecordNotFound untuk penanganan di Handler
		return "", nil, err
	}

	// 2. Verifikasi Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Password salah
		return "", nil, models.ErrInvalidCredentials
	}

	// 3. Buat JWT Token
	tokenString, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	// Sembunyikan password sebelum dikembalikan
	user.Password = ""
	return tokenString, user, nil
}
