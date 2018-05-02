package dnsmasq

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// NewLog Contructor
func NewLog() *Log {
	return new(Log)
}

// Log Data structure of DNSMASQ log line
type Log struct {
	gorm.Model
	Date     time.Time `gorm:"not null;index:idx_date_process"`
	Process  string    `gorm:"not null;index:idx_date_process"`
	PID      int       `gorm:"not null;"`
	QID      int       `gorm:"not null;index:idx_qid"`
	SourceIP string    `gorm:"not null;index:idx_source_ip"`
	QType    string    `gorm:"not null;index:idx_qtype"`
	Query    string    `gorm:"not null;"`
}

var logRegex = regexp.MustCompile(`^([A-Z][a-z]{2} [ 1-3][0-9] [ 0-2][0-9]:[0-6]?[0-9]:[0-6]?[0-9]) ([a-z-]+)\[([0-9]+)\]: (.*)$`)

// ParseLine Parse dnsmasq log line
func (l *Log) ParseLine(line string) *Log {
	matches := logRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		log.Fatalln("matches == 0 :", line)
	}

	l.Date, _ = time.Parse(time.Stamp, matches[1])
	l.Process = matches[2]
	l.PID, _ = strconv.Atoi(matches[3])

	switch l.Process {
	case "dnsmasq":
		l.parseQuery(matches[4])
	case "dnsmasq-dhcp":
		l.parseDhcpLog(matches[4])
	}

	return l
}

func (l *Log) parseQuery(s string) *Log {
	tokens, size := strings.Split(s, " "), len(s)

	if size >= 1 {
		l.QID, _ = strconv.Atoi(tokens[0])
	}
	if l.QID > 0 {
		if size >= 2 {
			l.SourceIP = strings.Split(tokens[1], "/")[0]
		}
		if size >= 3 {
			l.QType = tokens[2]
		}
		if size >= 4 {
			l.Query = strings.Join(tokens[3:], " ")
		}
	}

	return l
}

func (l *Log) parseDhcpLog(s string) *Log {

	l.Query = s

	return l
}
