package user

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/common"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type BusRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Bus, error)
}

type UseCase struct {
	busRepository BusRepository
}

func (u *UseCase) Login(ctx context.Context, data dto.DriverLoginRequest) (*dto.DriverLoginResponse, error) {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		return nil, err
	}

	bus, err := u.busRepository.FindByUsername(ctx, data.Username)
	if err != nil {
		log.Printf("error when finding bus by username: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(bus.Password),
		[]byte(data.Password),
	)
	if err != nil {
		log.Printf("wrong password, err: %v", err)
		return nil, err
	}

	token, err := common.NewJWT(*bus)
	if err != nil {
		log.Printf("error when creating jwt: %v", err)
		return nil, err
	}

	return &dto.DriverLoginResponse{
		ID:       bus.ID,
		Number:   bus.Number,
		Plate:    bus.Plate,
		Status:   bus.Status,
		Route:    bus.Route,
		IsActive: bus.IsActive,
		Token:    token,
	}, nil
}

func NewUseCase(busRepository BusRepository) *UseCase {
	return &UseCase{
		busRepository: busRepository,
	}
}
