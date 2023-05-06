package api

import (
	"fmt"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/bus"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/location"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/user"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Postgres will be the client for psql database
	Postgres *gorm.DB
)

// internal
var (
	grpcLocationHandler *location.Handler
	grpcUserHandler     *user.Handler
)

func InjectDependency() {
	// external dependencies
	addr := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	KafkaWriterLocation = &kafka.Writer{
		Addr:                   kafka.TCP(addr...),
		Topic:                  "location",
		Balancer:               &kafka.LeastBytes{},
		BatchTimeout:           5 * time.Millisecond,
		AllowAutoTopicCreation: true,
		Compression:            kafka.Snappy,
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to database, with error: %s", err.Error())
	}
	Postgres = db

	// internal package
	locationRepository := location.NewRepository(KafkaWriterLocation)
	locationUseCase := location.NewUseCase(locationRepository)
	grpcLocationHandler = location.NewHandler(locationUseCase)

	busRepository := bus.NewRepository(Postgres)

	userUseCase := user.NewUseCase(busRepository)
	grpcUserHandler = user.NewHandler(userUseCase)
}

func InitGRPCServer() {
	srv := grpc.NewServer()
	pb.RegisterLocationServer(srv, grpcLocationHandler)
	pb.RegisterUserServer(srv, grpcUserHandler)

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
