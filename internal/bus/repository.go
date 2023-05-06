package bus

import (
	"context"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*domain.Bus, error) {
	var bus domain.Bus
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&bus).Error
	if err != nil {
		return nil, err
	}

	return &bus, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
