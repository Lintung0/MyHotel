package routes

import (
	"backend/internal/config"
	"backend/internal/infra/http/handlers"
	"backend/internal/infra/http/routes/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	roomHandler *handlers.RoomHandler,
	bookingHandler *handlers.BookingHandler,
	reviewHandler *handlers.ReviewHandler,
	cfg *config.Config,
) {
	// Public Routes (Tanpa autentikasi)
	public := app.Group("/api")

	// Auth Routes
	auth := public.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Room Routes (Public - Lihat dan Cari)
	rooms := public.Group("/rooms")
	rooms.Get("", roomHandler.GetAllRooms)
	rooms.Get("/:id", roomHandler.GetRoomByID)
	rooms.Post("/available", roomHandler.GetAvailableRooms)

	// Review Routes (Public - Lihat)
	reviews := public.Group("/reviews")
	reviews.Get("/room/:roomId", reviewHandler.GetRoomReviews)
	reviews.Get("/:id", reviewHandler.GetReviewByID)

	// Protected Routes (Memerlukan autentikasi)
	protected := app.Group("/api", middleware.JWTMiddleware(cfg))

	// Member Routes
	member := protected.Group("/member")

	// Booking Routes (Member)
	bookings := member.Group("/bookings")
	bookings.Post("", bookingHandler.CreateBooking)
	bookings.Get("", bookingHandler.GetMyBookings)
	bookings.Delete("/:id", bookingHandler.CancelBooking)

	// Review Routes (Member)
	memberReviews := member.Group("/reviews")
	memberReviews.Post("", reviewHandler.CreateReview)

	// Admin Routes
	admin := protected.Group("/admin", middleware.RoleMiddleware("admin"))

	// Room Management Routes (Admin)
	adminRooms := admin.Group("/rooms")
	adminRooms.Post("", roomHandler.CreateRoom)
	adminRooms.Put("/:id", roomHandler.UpdateRoom)
	adminRooms.Delete("/:id", roomHandler.DeleteRoom)

	// Room Image Management Routes (Admin)
	adminRoomImages := admin.Group("/rooms/:id/images")
	adminRoomImages.Post("", roomHandler.AddRoomImage)
	adminRoomImages.Delete("/:imageId", roomHandler.DeleteRoomImage)

	// Booking Management Routes (Admin)
	adminBookings := admin.Group("/bookings")
	adminBookings.Get("", bookingHandler.GetAllBookings)
	adminBookings.Put("/:id/payment-status", bookingHandler.UpdatePaymentStatus)

	// Review Management Routes (Admin)
	adminReviews := admin.Group("/reviews")
	adminReviews.Delete("/:id", reviewHandler.DeleteReview)
}
