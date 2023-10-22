package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	inventoryDB, err := NewDB("inventory")
	if err != nil {
		panic(err)
	}
	defer inventoryDB.Close()

	reservationDB, err := NewDB("reservation")
	if err != nil {
		panic(err)
	}
	defer reservationDB.Close()

	inventoryRepository := NewInventoryRepository(inventoryDB)
	reservationRepository := NewReservationRepository(reservationDB)
	reservationService := NewService(inventoryRepository, reservationRepository)
	reservationHandler := NewHandler(reservationService)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	reservationHandler.RegisterRoutes(app)

	app.Listen(":3000")
}
