package api

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/bus"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/location"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/user"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"strings"
)

// GrpcSrv for serving grpc services
var GrpcSrv *grpc.Server

// External services
var (
	// Postgres will be the client for psql database
	Postgres *gorm.DB

	Producer sarama.SyncProducer
)

// internal
var (
	grpcLocationHandler *location.Handler
	grpcUserHandler     *user.Handler
)

func InjectDependency() {
	// external dependencies
	addr := strings.Split(os.Getenv("KAFKA_ADDR"), ";")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		log.Println("Failed to start Sarama producer:", err)
	}
	log.Printf("connect to kafka with broker: %v\n", addr)
	Producer = producer

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
	locationRepository := location.NewRepository(Producer)
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
