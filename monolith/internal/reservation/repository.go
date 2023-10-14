package reservation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Abdulsametileri/hotel-reservation-service/monolith/pkg/database"
)

var (
	ErrNoEnoughCapacity = errors.New("no enough capacity")
	ErrEditConflict     = errors.New("edit conflict")
)

type Repository interface {
	Create(ctx context.Context, reservation *Reservation) error
}

type defaultRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) Repository {
	return &defaultRepository{db: db}
}

func (repo *defaultRepository) Create(ctx context.Context, reservation *Reservation) error {
	tx, err := repo.db.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error initializing transaction %w", err)
	}
	defer tx.Rollback()

	inventory, err := checkInventory(tx, reservation)
	if err != nil {
		return fmt.Errorf("error checking inventory %w", err)
	}

	if inventory.TotalReserved+1 > inventory.TotalInventory {
		return ErrNoEnoughCapacity
	}

	if err = createReservation(tx, reservation); err != nil {
		return fmt.Errorf("error inserting reservation %w", err)
	}

	if err = updateInventory(tx, reservation, inventory.Version); err != nil {
		return fmt.Errorf("error updating inventory %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error during commiting %w", err)
	}

	return nil
}

func checkInventory(tx *sql.Tx, reservation *Reservation) (RoomTypeInventory, error) {
	checkInventoryQuery := `
		SELECT date, total_inventory, total_reserved, version
		FROM room_type_inventory
		WHERE hotel_id = $1 and room_type_id = $2 AND date between $3 AND $4`

	var roomTypeInventory RoomTypeInventory

	checkInventoryArgs := []interface{}{
		reservation.HotelID,
		reservation.RoomTypeID,
		reservation.StartDate,
		reservation.EndDate,
	}

	if err := tx.QueryRowContext(context.Background(), checkInventoryQuery, checkInventoryArgs...).
		Scan(
			&roomTypeInventory.Date,
			&roomTypeInventory.TotalInventory,
			&roomTypeInventory.TotalReserved,
			&roomTypeInventory.Version,
		); err != nil {
		return RoomTypeInventory{}, err
	}

	return roomTypeInventory, nil
}

func createReservation(tx *sql.Tx, reservation *Reservation) error {
	createReservationQuery := `
		INSERT INTO reservation (reservation_id, hotel_id, room_type_id, start_date, end_date, status) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	createReservationArgs := []interface{}{
		reservation.ReservationID,
		reservation.HotelID,
		reservation.RoomTypeID,
		reservation.StartDate,
		reservation.EndDate,
		"pending_pay",
	}

	if _, err := tx.Exec(createReservationQuery, createReservationArgs...); err != nil {
		return err
	}

	return nil
}

func updateInventory(tx *sql.Tx, reservation *Reservation, version int) error {
	updateInventoryQuery := `
		UPDATE room_type_inventory
		SET total_reserved = total_reserved + 1, version = version + 1 
		WHERE hotel_id = $1 AND room_type_id = $2 AND date between $3 AND $4 
	  	AND version = $5`

	updateInventoryArgs := []interface{}{
		reservation.HotelID,
		reservation.RoomTypeID,
		reservation.StartDate,
		reservation.EndDate,
		version,
	}

	result, err := tx.Exec(updateInventoryQuery, updateInventoryArgs...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}
