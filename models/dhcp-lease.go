package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type DhcpLease struct {
	Expiry    time.Time
	Mac       string
	IP        string
	Name      string
	ClientID  []string
	Tag       string
	LeaseTime string
	Ignore    bool
}

func NewDhcpLease() *DhcpLease {
	return &DhcpLease{
		Ignore: false,
	}
}

// ParseLease Parse DNSMASQ dhcp lease file format
// Line format : 0123456789 01:02:03:04:05:06 123.456.789.123 abcdef 01:02:03:04:05:06:07
// 1 - Time of expiry , in Epoch time
// 2 - MAC Address
// 3 - IP Address
// 4 - Computer Name
// 5 - Client ID : Computer's unique ID
func (l *DhcpLease) ParseLease(s string) error {
	tokens := strings.Split(s, " ")
	size := len(tokens)

	if size < 3 {
		return errors.New("Parsed string has bad format")
	}

	switch size {
	case 5:
		l.ClientID = []string{tokens[4]}
		fallthrough
	case 4:
		l.Name = tokens[3]
		fallthrough
	case 3:
		l.IP = tokens[2]
		l.Mac = tokens[1]
		epoch, _ := strconv.ParseInt(tokens[0], 10, 64)
		l.Expiry = time.Unix(epoch, 0)
		fallthrough
	default:
	}

	return nil
}

// ParseStaticLease Parse DNSMASQ static dhcp lease config format
// dhcp-host=[<hwaddr>][,id:<client_id>|*][,set:<tag>][,<ipaddr>][,<hostname>][,<lease_time>][,ignore]
func (l *DhcpLease) ParseStaticLease(s string) error {
	if !strings.HasPrefix(s, "dhcp-host=") {
		return errors.New("Parsed string is not dhcp-host string")
	}

	tokens := strings.Split(strings.Replace(s, "dhcp-host=", "", 1), ",")
	size, id, set := len(tokens), 0, 0

	for i := 1; i < 2; i++ {
		if strings.HasPrefix(tokens[i], "id:") {
			id++
			l.ClientID = strings.Split(strings.Replace(tokens[i], "id:", "", 1), "|")
		}
		if strings.HasPrefix(tokens[i], "set:") {
			set++
			l.Tag = strings.Replace(tokens[i], "set:", "", 1)
		}
	}
	switch size {
	case 5 + id + set:
		l.Ignore = true
		fallthrough
	case 4 + id + set:
		l.LeaseTime = tokens[3+id+set]
		fallthrough
	case 3 + id + set:
		l.Name = tokens[2+id+set]
		fallthrough
	case 2 + id + set:
		l.IP = tokens[1+id+set]
		fallthrough
	case 1:
		l.Mac = tokens[0]
		fallthrough
	default:
	}

	return nil
}
