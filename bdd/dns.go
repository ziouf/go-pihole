package bdd

import (
	"encoding/json"
	"time"
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

func (d *DNS) Encode() ([]byte, []byte) {
	key := []byte(d.Date.Format(time.Stamp))
	value, _ := json.Marshal(d)
	return key, value
}

func (d *DNS) Decode(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	return nil
}
