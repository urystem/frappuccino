package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]any

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}
