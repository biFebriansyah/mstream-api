package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB[T any] struct {
	Data T
}

func (j *JSONB[T]) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	return json.Unmarshal(b, &j.Data)
}

func (j *JSONB[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}
