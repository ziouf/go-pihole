package bdd

import (
	"encoding/json"
	"time"
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

func (d *DHCP) Encode() ([]byte, []byte) {
	key := []byte(d.Date.Format(time.Stamp))
	value, _ := json.Marshal(d)
	return key, value
}

func (d *DHCP) Decode(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	return nil
}
