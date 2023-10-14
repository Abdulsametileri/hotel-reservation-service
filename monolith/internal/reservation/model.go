package reservation

type Reservation struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	HotelID       int    `json:"hotelId"`
	RoomTypeID    int    `json:"roomTypeId"`
	ReservationID string `json:"reservationId"`
}

type RoomTypeInventory struct {
	HotelID        int    `json:"hotelId"`
	RoomTypeID     int    `json:"roomTypeId"`
	Date           string `json:"date"`
	TotalInventory int    `json:"totalInventory"`
	TotalReserved  int    `json:"totalReserved"`
	Version        int    `json:"-"`
}
