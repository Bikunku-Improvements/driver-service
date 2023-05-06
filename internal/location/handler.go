package location

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/common"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type useCase interface {
	SendLocation(ctx context.Context, loc domain.Location, bus domain.Bus) error
}

type Handler struct {
	pb.UnimplementedLocationServer
	useCase useCase
}

func (h Handler) SendLocation(server pb.Location_SendLocationServer) error {
	md, ok := metadata.FromIncomingContext(server.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "login required")
	}

	token := md.Get("token")
	if len(token) <= 0 {
		return status.Errorf(codes.Unauthenticated, "login required")
	}

	claims, err := common.ExtractTokenData(token[0])
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "token expired or invalid")
	}

	for {
		req, err := server.Recv()
		if err == io.EOF {
			// Close the connection and return the response to the client
			return server.SendAndClose(&pb.SendLocationResponse{Message: "OK"})
		}
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		loc := domain.Location{
			BusID:     claims.Data.ID,
			Long:      float64(req.GetLong()),
			Lat:       float64(req.GetLat()),
			Speed:     float64(req.GetSpeed()),
			Heading:   float64(req.GetHeading()),
			CreatedAt: time.Now(),
		}

		err = h.useCase.SendLocation(server.Context(), loc, claims.Data)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send location: "+err.Error())
		}
	}
}

func NewHandler(useCase useCase) *Handler {
	return &Handler{useCase: useCase}
}
