package reservation

type ReqBody struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	HotelID       int    `json:"hotelId"`
	RoomTypeID    int    `json:"roomTypeId"`
	ReservationID string `json:"reservationId"`
}

func (r *ReqBody) Validate() error {
	return nil
}

func (r *ReqBody) ToReservation() *Reservation {
	return &Reservation{
		StartDate:     r.StartDate,
		EndDate:       r.EndDate,
		HotelID:       r.HotelID,
		RoomTypeID:    r.RoomTypeID,
		ReservationID: r.ReservationID,
	}
}
