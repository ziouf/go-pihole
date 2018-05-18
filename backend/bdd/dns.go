package bdd

import (
	"encoding/json"
)

// DNS
type DNS struct {
	*Log

	QID     uint
	From    string
	Type    string
	FQDN    string
	IP      string
	Secured bool
}

func (d *DNS) Encode() []byte {
	value, _ := json.Marshal(d)
	return value
}

func (d *DNS) Decode(data []byte) error {
	return json.Unmarshal(data, d)
}
