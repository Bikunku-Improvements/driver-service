package location

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
)

type UseCase struct {
	repository repository
}

type repository interface {
	SendLocation(ctx context.Context, loc dto.SendLocationDataDTO) error
}

func (u UseCase) SendLocation(ctx context.Context, loc domain.Location, bus domain.Bus) error {
	return u.repository.SendLocation(ctx, dto.SendLocationDataDTO{
		BusID:     bus.ID,
		Number:    bus.Number,
		Plate:     bus.Plate,
		Status:    bus.Status,
		Route:     bus.Route,
		IsActive:  bus.IsActive,
		Long:      loc.Long,
		Lat:       loc.Lat,
		Speed:     loc.Speed,
		Heading:   loc.Heading,
		CreatedAt: loc.CreatedAt,
	})
}

func NewUseCase(repository repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}
