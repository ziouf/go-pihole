package parser

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cm-cloud.fr/go-pihole/bdd"
	"github.com/gobuffalo/uuid"
)

var logRegex = regexp.MustCompile(`^([A-Z][a-z]{2} [ 1-3][0-9] [ 0-2][0-9]:[0-6]?[0-9]:[0-6]?[0-9]) ([a-z-]+)\[([0-9]+)\]: (.*)$`)

// LogParser
var LogParser = new(logParse)

type logParse struct{}

func (lp *logParse) ParseLine(line string) bdd.Encodable {
	matches := logRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		log.Fatalln("matches == 0 :", line)
	}

	id, _ := uuid.NewV1()
	date, _ := time.Parse(time.Stamp, matches[1])
	pid, _ := strconv.Atoi(matches[3])

	l := &bdd.Log{
		Model:   bdd.Model{ID: id, CreatedAt: time.Now()},
		Date:    date,
		Process: matches[2],
		PID:     pid,
	}

	switch l.Process {
	case `dnsmasq`:
		return lp.parseDNSLine(l, matches[4])
	case `dnsmasq-dhcp`:
		return lp.parseDhcpLine(l, matches[4])
	// case `dnsmasq-script`:
	// case `dnsmasq-tftp`:
	default:
		return l
	}
}

// DNS query mapping
// QID  : Displayed in dnsmasq log if query as extra
// Type : query[A], query[AAAA], cached, forwarded, validation, reply
func (lp *logParse) parseDNSLine(log *bdd.Log, line string) *bdd.DNS {
	dns := &bdd.DNS{Log: log}
	tokens := strings.Split(line, " ")
	size := len(tokens)

	if size == 0 {
		return dns
	}
	if qid, err := strconv.Atoi(tokens[0]); err == nil {
		dns.QID = uint(qid)
	}
	if size > 1 {
		dns.From = strings.Split(tokens[1], "/")[0]
	}
	if size > 2 {
		dns.Type = tokens[2]
	}
	if size > 3 {
		switch dns.Type {
		case `query[A]`, `query[AAAA]`:
			dns.FQDN = tokens[3]
		case `cached`, `reply`:
			dns.FQDN = tokens[3]
			dns.IP = tokens[5]
		case `validation`:
			dns.Secured = tokens[5] != `INSECURE`
		default:
			dns.FQDN = tokens[3]
		}
	}

	return dns
}

func (lp *logParse) parseDhcpLine(log *bdd.Log, line string) *bdd.DHCP {
	dhcp := &bdd.DHCP{Log: log}
	tokens := strings.Split(line, " ")

	switch strings.Split(tokens[0], `(`)[0] {
	case `DHCPDISCOVER`, `RTR-SOLICIT`:
		dhcp.Type = tokens[0]
		dhcp.MAC = tokens[1]
	case `RTR-ADVERT`:
		dhcp.Type = tokens[0]
		dhcp.IP = tokens[1]
	case `DHCPOFFER`, `DHCPREQUEST`:
		dhcp.Type = tokens[0]
		dhcp.IP = tokens[1]
		dhcp.MAC = tokens[2]
	case `SLAAC-CONFIRM`:
		dhcp.Type = tokens[0]
		dhcp.IP = tokens[1]
		dhcp.Name = tokens[2]
	case `DHCPACK`:
		dhcp.Type = tokens[0]
		dhcp.IP = tokens[1]
		dhcp.MAC = tokens[2]
		dhcp.Name = tokens[3]
	default:
		dhcp.Msg = line
	}

	return dhcp
}
