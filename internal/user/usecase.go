package user

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/common"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/logger"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		logger.Logger.Error("error when finding bus by username", zap.Error(err), zap.Any("data", data))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(bus.Password),
		[]byte(data.Password),
	)
	if err != nil {
		logger.Logger.Error("wrong password", zap.Error(err))
		return nil, err
	}

	token, err := common.NewJWT(*bus)
	if err != nil {
		logger.Logger.Error("error when creating jwt", zap.Error(err), zap.Any("data", bus))
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
