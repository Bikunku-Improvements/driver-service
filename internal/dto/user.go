package dto

import "github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"

type DriverLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DriverLoginResponse struct {
	ID       uint             `json:"id"`
	Number   int              `json:"number"`
	Plate    string           `json:"plate"`
	Status   domain.BusStatus `json:"status"`
	Route    domain.Route     `json:"route"`
	IsActive bool             `json:"isActive"`
	Token    string           `json:"token"`
}
