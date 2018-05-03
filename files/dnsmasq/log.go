package dnsmasq

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NewLog Contructor
func NewLog() *Log {
	return new(Log)
}

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
}

// Log Data structure of DNSMASQ log line
type Log struct {
	Model
	Date    time.Time `gorm:"not null;index:idx_date_process;"`
	Process string    `gorm:"not null;index:idx_date_process;"`
	PID     int       `gorm:"not null;"`

	DnsQID     int
	DnsFrom    string
	DnsType    string
	DnsFQDN    string
	DnsIP      string
	DnsSecured bool

	DhcpType string
	DhcpIP   string
	DhcpMAC  string
	DhcpName string
	DhcpMsg  string
}

// DNS query mapping
// QID  : Displayed in dnsmasq log if query as extra
// Type : query[A], query[AAAA], cached, forwarded, validation, reply

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
	case `dnsmasq`:
		l.parseDNSLog(matches[4])
	case `dnsmasq-dhcp`:
		l.parseDhcpLog(matches[4])
	case `dnsmasq-script`:
		// TODO : implement parsing script
	case `dnsmasq-tftp`:
		// TODO : implement parsing script
	}

	return l
}

func (l *Log) parseDNSLog(s string) *Log {
	tokens := strings.Split(s, " ")
	size := len(tokens)
	
	if size >= 1 {
		l.DnsQID, _ = strconv.Atoi(tokens[0])
	}
	if l.DnsQID > 0 {
		if size >= 2 {
			l.DnsFrom = strings.Split(tokens[1], "/")[0]
		}
		if size >= 3 {
			l.DnsType = tokens[2]
		}
		if size >= 4 {
			switch l.DnsType {
			case `query[A]`, `query[AAAA]`:
				l.DnsFQDN = tokens[3]
			case `cached`, `reply`:
				l.DnsFQDN = tokens[3]
				l.DnsIP = tokens[5]
			case `validation`:
				l.DnsSecured = tokens[5] != `INSECURE`
			default:
				l.DnsFQDN = tokens[3]
			}
		}
	}

	return l
}

func (l *Log) parseDhcpLog(s string) *Log {
	tokens := strings.Split(s, " ")

	switch strings.Split(tokens[0], `(`)[0] {
	case `DHCPDISCOVER`, `RTR-SOLICIT`:
		l.DhcpType = tokens[0]
		l.DhcpMAC = tokens[1]
	case `RTR-ADVERT`:
		l.DhcpType = tokens[0]
		l.DhcpIP = tokens[1]
	case `DHCPOFFER`, `DHCPREQUEST`:
		l.DhcpType = tokens[0]
		l.DhcpIP = tokens[1]
		l.DhcpMAC = tokens[2]
	case `SLAAC-CONFIRM`:
		l.DhcpType = tokens[0]
		l.DhcpIP = tokens[1]
		l.DhcpName = tokens[2]
	case `DHCPACK`:
		l.DhcpType = tokens[0]
		l.DhcpIP = tokens[1]
		l.DhcpMAC = tokens[2]
		l.DhcpName = tokens[3]
	default:
		l.DhcpMsg = s
	}

	return l
}
