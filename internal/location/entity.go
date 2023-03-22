package location

import (
	"encoding/json"
	"time"
)

type Location struct {
	Long      string    `json:"long"`
	Lat       string    `json:"lat"`
	CreatedAt time.Time `json:"created_at"`
	BusID     string    `json:"bus_id"`
}

func (l Location) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l Location) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &l); err != nil {
		return err
	}

	return nil
}
