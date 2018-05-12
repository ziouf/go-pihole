package bdd

import (
	"encoding/json"
)

// DHCP
type DHCP struct {
	*Log

	Type string
	IP   string
	MAC  string
	Name string
	Msg  string
}

func (d *DHCP) Encode() []byte {
	value, _ := json.Marshal(d)
	return value
}

func (d *DHCP) Decode(data []byte) error {
	return json.Unmarshal(data, d)
}
