package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const (
	EndpointCreateReservation = "/v1/reservations"
	EndpointGetReservation    = "/v1/reservation/:id"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get(EndpointGetReservation, h.GetReservation)
	app.Post(EndpointCreateReservation, h.CreateReservation)
}

func (h *Handler) GetReservation(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotImplemented)
}

func (h *Handler) CreateReservation(ctx *fiber.Ctx) error {
	var reqBody ReqBody
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := reqBody.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := h.service.CreateReservation(ctx.Context(), reqBody.ToReservation()); err != nil {
		switch {
		case errors.Is(err, ErrNoEnoughCapacity):
			return ctx.Status(http.StatusBadRequest).SendString("Not enough capacity")
		default:
			fmt.Println("[ERROR]", err.Error())
			return ctx.Status(http.StatusInternalServerError).SendString("Please try again later")
		}

	}

	ctx.Location("/v1/reservation/" + reqBody.ReservationID)
	return ctx.SendStatus(http.StatusCreated)
}
