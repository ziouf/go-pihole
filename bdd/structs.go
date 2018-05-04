package bdd

import (
	"encoding/json"
	"time"
)

type Model struct {
	Stamp time.Time
}

func (m *Model) Encode() ([]byte, error) {
	return json.Marshal(m)
}
func (m *Model) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m *Model) StampEncoded() []byte {
	return []byte(m.Stamp.Format(time.RFC3339))
}
