package domain

import (
	"encoding/json"
	"time"
)

// Crowded Status
const (
	EMPTY    BusStatus = "EMPTY"
	MODERATE BusStatus = "MODERATE"
	FULL     BusStatus = "FULL"
)

// Route Type
const (
	RED  Route = "RED"
	BLUE Route = "BLUE"
)

type BusStatus string

type Route string

type Bus struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Number   int       `gorm:"column:number;unique" json:"number"`
	Plate    string    `gorm:"column:plate;unique" json:"plate"`
	Status   BusStatus `gorm:"column:status;default:EMPTY" json:"status"`
	Route    Route     `gorm:"column:route" json:"route"`
	IsActive bool      `gorm:"column:is_active;default:false" json:"is_active"`
	Username string    `gorm:"column:username;unique" json:"username"`
	Password string    `gorm:"password" json:"-"`
}

func (b *Bus) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Bus) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, b); err != nil {
		return err
	}

	return nil
}

type Location struct {
	BusID     uint      `json:"bus_id"`
	Long      float64   `json:"long"`
	Lat       float64   `json:"lat"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Location) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *Location) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, l); err != nil {
		return err
	}

	return nil
}
