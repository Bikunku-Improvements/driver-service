package location

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type useCase interface {
	SendLocation(ctx context.Context, loc Location) error
}

type Handler struct {
	pb.UnimplementedLocationServer
	useCase useCase
}

func (h Handler) SendLocation(server pb.Location_SendLocationServer) error {
	for {
		req, err := server.Recv()
		if err == io.EOF {
			// Close the connection and return the response to the client
			return server.SendAndClose(&pb.SendLocationResponse{Message: "OK"})
		}
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}

		loc := Location{
			Long:      req.GetLong(),
			Lat:       req.GetLat(),
			BusID:     req.GetBusId(),
			CreatedAt: time.Now(),
		}

		err = h.useCase.SendLocation(server.Context(), loc)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
}

func NewHandler(useCase useCase) *Handler {
	return &Handler{useCase: useCase}
}
