package api

import (
	"fmt"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/location"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// GrpcSrv for serving grpc services
var GrpcSrv *grpc.Server

// External services
var (
	// KafkaWriterLocation will be producer of location
	KafkaWriterLocation *kafka.Writer
)

// internal
var (
	grpcLocationHandler *location.Handler
)

func InjectDependency() {
	// external dependencies
	addr := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	KafkaWriterLocation = &kafka.Writer{
		Addr:         kafka.TCP(addr...),
		Topic:        "location",
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchTimeout: 10 * time.Millisecond,
	}

	// internal package
	locationRepository := location.NewRepository(KafkaWriterLocation)
	locationUseCase := location.NewUseCase(locationRepository)
	grpcLocationHandler = location.NewHandler(locationUseCase)
}

func InitGRPCServer() {
	srv := grpc.NewServer()
	pb.RegisterLocationServer(srv, grpcLocationHandler)

	GrpcSrv = srv
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := GrpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}