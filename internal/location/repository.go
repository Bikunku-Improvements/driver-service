package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/dto"
	"log"
)

type Repository struct {
	producer sarama.SyncProducer
}

//func (r Repository) SendLocation(ctx context.Context, loc dto.SendLocationDataDTO) error {
//	b, err := json.Marshal(loc)
//	if err != nil {
//		return fmt.Errorf("failed to unmarshal location: %v", err)
//	}
//
//	msg := kafka.Message{
//		Value: b,
//	}
//
//	err = r.writer.WriteMessages(ctx, msg)
//	if err != nil {
//		return fmt.Errorf("failed to write message: %v", err)
//	}
//	log.Printf("message send to %v topic with value %v", r.writer.Topic, string(msg.Value))
//
//	return nil
//}

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
	log.Printf("message send to %v topic with value %v", "location", string(msgBytes))
	return nil
}

func NewRepository(producer sarama.SyncProducer) *Repository {
	return &Repository{
		producer: producer,
	}
}
