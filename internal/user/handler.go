package user

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
)

type useCase interface {
	Login(ctx context.Context, data dto.DriverLoginRequest) (*dto.DriverLoginResponse, error)
}

type Handler struct {
	pb.UnimplementedUserServer
	useCase useCase
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	payload := dto.DriverLoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	resp, err := h.useCase.Login(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Id:       uint64(resp.ID),
		Number:   int64(resp.Number),
		Plate:    resp.Plate,
		Status:   string(resp.Status),
		Route:    string(resp.Route),
		IsActive: resp.IsActive,
		Token:    resp.Token,
	}, nil
}

func NewHandler(useCase useCase) *Handler {
	return &Handler{useCase: useCase}
}
