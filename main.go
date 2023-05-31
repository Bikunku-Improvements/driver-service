package main

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/api"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger.Logger, _ = config.Build()

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Error("error when loading .env", zap.Error(err))
	}
	go func() {
		api.InjectDependency()
		api.InitGRPCServer()
		api.StartGRPCServer()
	}()

	<-ctx.Done()
	stop()
	api.GrpcSrv.GracefulStop()
	api.Producer.Close()
	dbInstance, _ := api.Postgres.DB()
	dbInstance.Close()
	logger.Logger.Sync()
}
