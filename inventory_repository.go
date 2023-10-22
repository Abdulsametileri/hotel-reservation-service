package main

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNoEnoughCapacity = errors.New("no enough capacity")
)

type InventoryRepository interface {
	CommitPrepared(ctx context.Context, txID string)
	RollbackPrepared(ctx context.Context, txID string)
	UpdatePrepared(ctx context.Context, r *Reservation, txID string) (RoomTypeInventory, error)
}

type defaultInventoryRepository struct {
	db *DB
}

func NewInventoryRepository(db *DB) InventoryRepository {
	return &defaultInventoryRepository{db: db}
}

func (repo *defaultInventoryRepository) UpdatePrepared(ctx context.Context, r *Reservation, txID string) (RoomTypeInventory, error) {
	tx, err := repo.db.DB().BeginTx(ctx, nil)
	if err != nil {
		return RoomTypeInventory{}, fmt.Errorf("error begin statement %w", err)
	}

	selectQuery := `
		SELECT date, total_inventory, total_reserved
		FROM room_type_inventory
		WHERE hotel_id = $1 AND room_type_id = $2 AND date BETWEEN $3 AND $4
		FOR UPDATE;
	`

	args := []interface{}{
		r.HotelID,
		r.RoomTypeID,
		r.StartDate,
		r.EndDate,
	}

	var inventory RoomTypeInventory
	if err = tx.QueryRowContext(ctx, selectQuery, args...).
		Scan(&inventory.Date, &inventory.TotalInventory, &inventory.TotalReserved); err != nil {
		return RoomTypeInventory{}, err
	}

	updateQuery := `
		UPDATE room_type_inventory
		SET total_reserved = total_reserved + 1
		WHERE hotel_id = $1 AND room_type_id = $2 AND date BETWEEN $3 AND $4
	`

	if _, err = tx.ExecContext(ctx, updateQuery, args...); err != nil {
		return RoomTypeInventory{}, err
	}

	if _, err = tx.ExecContext(ctx, fmt.Sprintf("PREPARE TRANSACTION '%s';", txID)); err != nil {
		return RoomTypeInventory{}, err
	}

	if inventory.TotalReserved > inventory.TotalInventory {
		return RoomTypeInventory{}, ErrNoEnoughCapacity
	}

	return inventory, nil
}

func (repo *defaultInventoryRepository) CommitPrepared(ctx context.Context, txID string) {
	repo.db.DB().ExecContext(ctx, fmt.Sprintf("COMMIT PREPARED '%s';", txID))
}

func (repo *defaultInventoryRepository) RollbackPrepared(ctx context.Context, txID string) {
	repo.db.DB().ExecContext(ctx, fmt.Sprintf("ROLLBACK PREPARED '%s';", txID))
}
