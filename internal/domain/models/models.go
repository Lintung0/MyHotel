package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// --users--

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	FullName string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"type:enum('admin', 'member');default:'member'"`

	// Relasi: User punya banyak Booking
	Bookings []Booking `gorm:"foreignKey:UserID"`
}

type Room struct {
	gorm.Model
	RoomNumber   string  `gorm:"type:varchar(10);unique;not null"`
	Type         string  `gorm:"type:varchar(50);not null"`
	Price        float64 `gorm:"type:decimal(10,2);not null"`
	Description  string  `gorm:"type:text"`
	Status       string  `gorm:"type:enum('available', 'booked', 'maintenance');default:'available'"`
	MaxOccupancy int     `gorm:"not null"`

	// Relasi: Room punya banyak Image dan Booking
	Images   []RoomImage `gorm:"foreignKey:RoomID"`
	Bookings []Booking   `gorm:"foreignKey:RoomID"`
}

type RoomImage struct {
	gorm.Model
	RoomID    uint   `gorm:"not null"` // Foreign Key
	ImageURL  string `gorm:"type:varchar(255);not null"`
	IsPrimary bool   `gorm:"default:false"`
}

type Booking struct {
	gorm.Model
	UserID        uint      `gorm:"not null"` // Foreign Key ke User
	RoomID        uint      `gorm:"not null"` // Foreign Key ke Room
	CheckInDate   time.Time `gorm:"type:date;not null"`
	CheckOutDate  time.Time `gorm:"type:date;not null"`
	TotalPrice    float64   `gorm:"type:decimal(10,2);not null"`
	PaymentMethod string    `gorm:"type:varchar(50)"`
	PaymentStatus string    `gorm:"type:enum('pending', 'paid', 'failed');default:'pending'"`
	BookingStatus string    `gorm:"type:enum('confirmed', 'cancelled', 'completed');default:'confirmed'"`

	// Relasi: Booking punya 1 Review
	Review Review `gorm:"foreignKey:BookingID"`
}

type Review struct {
	gorm.Model
	BookingID uint   `gorm:"unique;not null"` // Foreign Key ke Booking (Unique)
	UserID    uint   `gorm:"not null"`        // Untuk kemudahan query
	Rating    int    `gorm:"type:int;not null;check:rating >= 1 AND rating <= 5"`
	Comment   string `gorm:"type:text"`
}

// -- pagination --
type Pagination struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Sort   string `json:"sort"` // Contoh: "created_at desc"
	Offset int    `json:"-"`
}

// --- JWT Claims ---
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// --- Konstanta Role ---
const (
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// --- Status Pembayaran ---
const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusFailed  = "failed"
)

// --- Status Pemesanan ---
const (
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

// --- Custom Errors ---
var (
	ErrRecordNotFound     = gorm.ErrRecordNotFound
	ErrInvalidCredentials = errors.New("username atau password salah")
	// Tambahkan error lain sesuai kebutuhan (misalnya: errors.New("kamar sudah dibooking"))
)
