package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/logger"
	"go.uber.org/zap"
)

type Repository struct {
	producer sarama.SyncProducer
}

func (r Repository) SendLocationWithSarama(ctx context.Context, loc dto.SendLocationDataDTO) error {
	msgBytes, err := json.Marshal(loc)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	_, _, err = r.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "location",
		Value: sarama.ByteEncoder(msgBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	logger.Logger.Info("message send to kafka", zap.String("topic", "location"), zap.String("value", string(msgBytes)))
	return nil
}

func NewRepository(producer sarama.SyncProducer) *Repository {
	return &Repository{
		producer: producer,
	}
}
