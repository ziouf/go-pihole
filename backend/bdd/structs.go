package bdd

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (m *Model) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Model) Encode() []byte {
	enc, _ := json.Marshal(m)
	return enc
}
