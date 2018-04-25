package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DnsmasqLog Data structure of DNSMASQ log line
type DnsmasqLog struct {
	gorm.Model
	Date     time.Time
	Process  string
	PID      int
	QID      int
	SourceIP string
	QType    string
	Query    string
}
