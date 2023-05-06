package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
	"github.com/segmentio/kafka-go"
	"log"
)

type Repository struct {
	writer *kafka.Writer
}

func (r Repository) SendLocation(ctx context.Context, loc dto.SendLocationDataDTO) error {
	b, err := json.Marshal(loc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal location: %v", err)
	}

	msg := kafka.Message{
		Value: b,
	}

	err = r.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}
	log.Printf("message send to %v topic with value %v", r.writer.Topic, string(msg.Value))

	return nil
}

func NewRepository(writer *kafka.Writer) *Repository {
	return &Repository{
		writer: writer,
	}
}
