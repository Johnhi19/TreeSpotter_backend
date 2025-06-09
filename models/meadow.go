package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IntSlize []int

type Meadow struct {
	ID       int      `json:"id"`
	Location string   `json:"location"`
	Name     string   `json:"name"`
	Size     IntSlize `json:"size"`
	TreeIds  IntSlize `json:"treeIds"`
}

func (s *IntSlize) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert DB value to []byte")
	}
	return json.Unmarshal(bytes, s)
}

func (s IntSlize) Value() (driver.Value, error) {
	return json.Marshal(s)
}
