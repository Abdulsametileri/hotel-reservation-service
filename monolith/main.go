package main

import (
	"github.com/Abdulsametileri/hotel-reservation-service/monolith/internal/reservation"
	"github.com/Abdulsametileri/hotel-reservation-service/monolith/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	reservationRepository := reservation.NewRepository(db)
	reservationService := reservation.NewService(reservationRepository)
	reservationHandler := reservation.NewHandler(reservationService)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	reservationHandler.RegisterRoutes(app)

	app.Listen(":3000")
}
