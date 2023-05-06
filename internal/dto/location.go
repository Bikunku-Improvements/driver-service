package dto

import (
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"time"
)

type SendLocationDataDTO struct {
	BusID     uint             `json:"bus_id"`
	Number    int              `json:"number"`
	Plate     string           `json:"plate"`
	Status    domain.BusStatus `json:"status"`
	Route     domain.Route     `json:"route"`
	IsActive  bool             `json:"isActive"`
	Long      float64          `json:"long"`
	Lat       float64          `json:"lat"`
	Speed     float64          `json:"speed"`
	Heading   float64          `json:"heading"`
	CreatedAt time.Time        `json:"created_at"`
}
