package main

import (
	"backend/internal/app/services"
	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/infra/database/mysql"
	"backend/internal/infra/gorm/repositories"
	"backend/internal/infra/http/handlers"
	"backend/internal/infra/http/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Database
	db := mysql.InitDB(cfg)

	// 3. Auto Migrate All Models
	mysql.AutoMigrate(db,
		&models.User{},
		&models.Room{},
		&models.RoomImage{},
		&models.Booking{},
		&models.Review{},
	)

	// 4. Initialize Repositories
	userRepo := repositories.NewGormRepository(db)
	roomRepo := repositories.NewGormRoomRepository(db)
	bookingRepo := repositories.NewGormBookingRepository(db)
	roomImageRepo := repositories.NewGormRoomImageRepository(db)
	reviewRepo := repositories.NewGormReviewRepository(db)

	// 5. Initialize Services
	authService := services.NewAuthService(userRepo, cfg)
	roomService := services.NewRoomService(roomRepo, roomImageRepo)
	bookingService := services.NewBookingService(bookingRepo, roomRepo, reviewRepo)
	reviewService := services.NewReviewService(reviewRepo, bookingRepo)

	// 6. Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	roomHandler := handlers.NewRoomHandler(roomService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	// 7. Create Fiber App
	app := fiber.New()

	// 8. Add Middleware
	app.Use(logger.New())

	// 9. Setup Routes
	routes.SetupRoutes(app, authHandler, roomHandler, bookingHandler, reviewHandler, cfg)

	// 10. Start Server
	port := ":" + cfg.ServerPort
	log.Printf("🚀 Server berjalan di http://localhost%s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
