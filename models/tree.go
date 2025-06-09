package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Tree struct {
	ID        int       `json:"id"`
	PlantDate time.Time `json:"plantDate"`
	MeadowId  int       `json:"meadowId"`
	Position  Position  `json:"position"`
	Type      string    `json:"type"`
}

func (p *Position) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert position to []byte")
	}

	if err := json.Unmarshal(bytes, p); err != nil {
		return fmt.Errorf("failed to unmarshal position: %w", err)
	}
	return nil
}
