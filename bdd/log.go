package bdd

import (
	"encoding/json"
	"time"
)

// Log
type Log struct {
	Model

	Date    time.Time
	Process string
	PID     int
}

func (d *Log) Encode() ([]byte, []byte) {
	key := []byte(d.Date.Format(time.Stamp))
	value, _ := json.Marshal(d)
	return key, value
}

func decode(name string, b []byte) *Log {
	switch name {
	case `DNS`:
		var dns DNS
		json.Unmarshal(b, &dns)
		return dns.Log
	case `DHCP`:
		var dhcp DHCP
		json.Unmarshal(b, &dhcp)
		return dhcp.Log
	default:
		return nil
	}
}
