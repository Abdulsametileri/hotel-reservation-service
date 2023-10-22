package main

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	CreateReservation(context.Context, *Reservation) error
}

type defaultService struct {
	inventoryRepository   InventoryRepository
	reservationRepository ReservationRepository
}

func NewService(inventoryRepository InventoryRepository, reservationRepository ReservationRepository) Service {
	return &defaultService{inventoryRepository: inventoryRepository, reservationRepository: reservationRepository}
}

func (d *defaultService) CreateReservation(ctx context.Context, reservation *Reservation) error {
	inventoryTxID := uuid.NewString()
	reservationTxID := uuid.NewString()

	_, err := d.inventoryRepository.UpdatePrepared(ctx, reservation, inventoryTxID)
	if err != nil {
		d.inventoryRepository.RollbackPrepared(ctx, inventoryTxID)
		return err
	}

	if err = d.reservationRepository.CreatePrepared(ctx, reservation, reservationTxID); err != nil {
		d.inventoryRepository.RollbackPrepared(ctx, inventoryTxID)
		return err
	}

	d.inventoryRepository.CommitPrepared(ctx, inventoryTxID)
	d.reservationRepository.CommitPrepared(ctx, reservationTxID)

	return nil
}
