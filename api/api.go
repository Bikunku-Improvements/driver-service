package api

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/grpc/pb"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/bus"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/location"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/logger"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	config.Version = sarama.V1_1_0_0

	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	config.Net.SASL.Enable = true
	config.Net.SASL.User = os.Getenv("KAFKA_USERNAME")
	config.Net.SASL.Password = os.Getenv("KAFKA_PASSWORD")
	config.Net.SASL.Handshake = true
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256

	tlsConfig := tls.Config{}
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tlsConfig

	producer, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		logger.Logger.Fatal("failed to start producer", zap.Error(err))
	}
	logger.Logger.Info("connect to producer", zap.Any("broker", addr))
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
		logger.Logger.Fatal("failed to connect to database", zap.Error(err))
	}
	logger.Logger.Info("connect to database")
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
		logger.Logger.Fatal("failed to listen tcp", zap.Error(err))
	}

	logger.Logger.Info("server listening", zap.Any("port", lis.Addr()))
	if err = GrpcSrv.Serve(lis); err != nil {
		logger.Logger.Fatal("failed to server grpc", zap.Error(err))
	}

}
