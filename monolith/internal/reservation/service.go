package reservation

import "context"

type Service interface {
	CreateReservation(context.Context, *Reservation) error
}

type defaultService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &defaultService{repository: repository}
}

func (d *defaultService) CreateReservation(ctx context.Context, reservation *Reservation) error {
	return d.repository.Create(ctx, reservation)
}
