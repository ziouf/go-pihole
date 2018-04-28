package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// NewDnsmasqLog Contructor
func NewDnsmasqLog() *DnsmasqLog {
	return &DnsmasqLog{}
}

// DnsmasqLog Data structure of DNSMASQ log line
type DnsmasqLog struct {
	// gorm.Model
	Date     time.Time `gorm:"not null;index:date_process"`
	Process  string    `gorm:"not null;index:date_process"`
	PID      int       `gorm:"not null;index:date_process"`
	QID      int       `gorm:"not null:index:qid"`
	SourceIP string    `gorm:"not null"`
	QType    string    `gorm:"not null"`
	Query    string    `gorm:"not null"`
}

func (l *DnsmasqLog) String() string {
	return fmt.Sprintf("{Date: '%s', Process: '%s', PID: '%d', QID: '%d', SourceIP: '%s', Qtype: '%s', Query: '%s'}", l.Date, l.Process, l.PID, l.QID, l.SourceIP, l.QType, l.Query)
}

func (l *DnsmasqLog) ParseLine(line string) *DnsmasqLog {
	runes := []rune(line)
	dateStr, rightStr := strings.TrimSpace(string(runes[:16])), string(runes[16:])

	if date, err := time.Parse(time.Stamp, dateStr); err == nil {
		l.Date = date
	} else {
		log.Println(err)
	}

	tokens, size := strings.Split(strings.TrimSpace(rightStr), " "), len(rightStr)
	if size >= 0 {
		t := strings.Split(tokens[0], "[")
		l.Process = t[0]
		pid, _ := strconv.Atoi(string([]rune(t[1])[:len(t[1])-2]))
		l.PID = pid
	}
	if size >= 1 {
		qid, _ := strconv.Atoi(tokens[1])
		l.QID = qid
	}
	if l.QID > 0 {
		if size >= 2 {
			l.SourceIP = strings.Split(tokens[2], "/")[0]
		}
		if size >= 3 {
			l.QType = tokens[3]
		}
		if size >= 4 {
			l.Query = strings.Join(tokens[4:], " ")
		}
	}

	return l
}
