package main

import (
	"context"
	"fmt"
)

type ReservationRepository interface {
	CommitPrepared(ctx context.Context, txID string)
	RollbackPrepared(ctx context.Context, txID string)
	CreatePrepared(ctx context.Context, r *Reservation, txID string) error
}

type defaultReservationRepository struct {
	db *DB
}

func NewReservationRepository(db *DB) ReservationRepository {
	return &defaultReservationRepository{db: db}
}

func (repo *defaultReservationRepository) CreatePrepared(ctx context.Context, r *Reservation, txID string) error {
	tx, err := repo.db.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error begin statement %w", err)
	}

	query := `
		INSERT INTO reservation (reservation_id, hotel_id, room_type_id, start_date, end_date, status) 
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	args := []interface{}{
		r.ReservationID,
		r.HotelID,
		r.RoomTypeID,
		r.StartDate,
		r.EndDate,
		"pending_pay",
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, fmt.Sprintf("PREPARE TRANSACTION '%s';", txID)); err != nil {
		return err
	}

	return nil
}

func (repo *defaultReservationRepository) CommitPrepared(ctx context.Context, txID string) {
	repo.db.DB().ExecContext(ctx, fmt.Sprintf("COMMIT PREPARED '%s';", txID))
}

func (repo *defaultReservationRepository) RollbackPrepared(ctx context.Context, txID string) {
	repo.db.DB().ExecContext(ctx, fmt.Sprintf("ROLLBACK PREPARED '%s';", txID))
}
