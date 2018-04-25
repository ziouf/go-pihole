package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DnsmasqLog Data structure of DNSMASQ log line
type DnsmasqLog struct {
	gorm.Model
	Date     time.Time `gorm:not null;index:date_process`
	Process  string    `gorm:not null;index:date_process`
	PID      int       `gorm:not null`
	QID      int       `gorm:not null`
	SourceIP string    `gorm:not null`
	QType    string    `gorm:not null`
	Query    string    `gorm:not null`
}
