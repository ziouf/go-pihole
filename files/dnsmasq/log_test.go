package dnsmasq

import (
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	t.Run("Parse Dnsmasq dns log", testParseDnsmasqLogLine)
	t.Run("Parse Dnsmasq dhcp log", testParseDnsmasqDhcpLogLine)
}

func testParseDnsmasqLogLine(t *testing.T) {
	tcs := []string{
		`Apr  9 00:50:24 dnsmasq[755]: 78348 192.168.1.254/52101 query[A] a.root-servers.net from 192.168.1.254`,
		`Apr  9 00:52:38 dnsmasq[755]: 78399 192.168.1.200/41427 cached cloudsync-tw.synology.com is <CNAME>`,
		`Apr  9 00:52:36 dnsmasq[755]: 78395 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35/32969 query[AAAA] ads.nexage.com from 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`,
		`Apr  9 00:52:36 dnsmasq[755]: 78395 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35/32969 /etc/pihole/gravity.list ads.nexage.com is 2a01:cb15:8046:3200:5ddd:efb1:d506:7ed6`,
	}
	rr := []Log{
		{Date: time.Date(int(0), time.April, int(9), int(0), int(50), int(24), int(0), time.UTC),
			Process: `dnsmasq`, PID: 755, QID: 78348, SourceIP: `192.168.1.254`, QType: `query[A]`, Query: `a.root-servers.net from 192.168.1.254`},
		{Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(38), int(0), time.UTC),
			Process: `dnsmasq`, PID: 755, QID: 78399, SourceIP: `192.168.1.200`, QType: `cached`, Query: `cloudsync-tw.synology.com is <CNAME>`},
		{Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(36), int(0), time.UTC),
			Process: `dnsmasq`, PID: 755, QID: 78395, SourceIP: `2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`, QType: `query[AAAA]`, Query: `ads.nexage.com from 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`},
		{Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(36), int(0), time.UTC),
			Process: `dnsmasq`, PID: 755, QID: 78395, SourceIP: `2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`, QType: `/etc/pihole/gravity.list`, Query: `ads.nexage.com is 2a01:cb15:8046:3200:5ddd:efb1:d506:7ed6`},
	}

	for i, tc := range tcs {
		l := NewLog().ParseLine(tc)

		if l.Date != rr[i].Date {
			t.Errorf("TC[%d] Failed to parse date : Wanted [%s] but found [%s]", i, rr[i].Date, l.Date)
		}
		if l.Process != rr[i].Process {
			t.Errorf("TC[%d] Failed to parse pocess : Wanted [%s] but found [%s]", i, rr[i].Process, l.Process)
		}
		if l.PID != rr[i].PID {
			t.Errorf("TC[%d] Failed to parse pid : Wanted [%d] but found [%d]", i, rr[i].PID, l.PID)
		}
		if l.QID != rr[i].QID {
			t.Errorf("TC[%d] Failed to parse qid : Wanted [%d] but found [%d]", i, rr[i].QID, l.QID)
		}
		if l.SourceIP != rr[i].SourceIP {
			t.Errorf("TC[%d] Failed to parse source ip : Wanted [%s] but found [%s]", i, rr[i].SourceIP, l.SourceIP)
		}
		if l.QType != rr[i].QType {
			t.Errorf("TC[%d] Failed to parse qtype : Wanted [%s] but found [%s]", i, rr[i].QType, l.QType)
		}
		if l.Query != rr[i].Query {
			t.Errorf("TC[%d] Failed to parse query : Wanted [%s] but found [%s]", i, rr[i].Query, l.Query)
		}
	}
}

func testParseDnsmasqDhcpLogLine(t *testing.T) {
	tcs := []string{
		`Apr  9 00:52:33 dnsmasq-dhcp[755]: DHCPDISCOVER(eth0) cc:fb:65:63:4d:f9`,
		`Apr  9 00:52:33 dnsmasq-dhcp[755]: DHCPOFFER(eth0) 192.168.1.4 cc:fb:65:63:4d:f9`,
		`Apr  9 00:52:33 dnsmasq-dhcp[755]: L'option dhcp-option redondante 3 sera ignorée`,
	}
	rr := []Log{
		{Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(33), int(0), time.UTC),
			Process: `dnsmasq-dhcp`, PID: 755, QID: 0, SourceIP: ``, QType: ``, Query: `DHCPDISCOVER(eth0) cc:fb:65:63:4d:f9`},
		{Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(33), int(0), time.UTC),
			Process: `dnsmasq-dhcp`, PID: 755, QID: 0, SourceIP: ``, QType: ``, Query: `DHCPOFFER(eth0) 192.168.1.4 cc:fb:65:63:4d:f9`},
		{Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(33), int(0), time.UTC),
			Process: `dnsmasq-dhcp`, PID: 755, QID: 0, SourceIP: ``, QType: ``, Query: `L'option dhcp-option redondante 3 sera ignorée`},
	}

	for i, tc := range tcs {
		l := NewLog().ParseLine(tc)

		if l.Date != rr[i].Date {
			t.Errorf("TC[%d] Failed to parse date : Wanted [%s] but found [%s]", i, rr[i].Date, l.Date)
		}
		if l.Process != rr[i].Process {
			t.Errorf("TC[%d] Failed to parse pocess : Wanted [%s] but found [%s]", i, rr[i].Process, l.Process)
		}
		if l.PID != rr[i].PID {
			t.Errorf("TC[%d] Failed to parse pid : Wanted [%d] but found [%d]", i, rr[i].PID, l.PID)
		}
		if l.QID != rr[i].QID {
			t.Errorf("TC[%d] Failed to parse qid : Wanted [%d] but found [%d]", i, rr[i].QID, l.QID)
		}
		if l.SourceIP != rr[i].SourceIP {
			t.Errorf("TC[%d] Failed to parse source ip : Wanted [%s] but found [%s]", i, rr[i].SourceIP, l.SourceIP)
		}
		if l.QType != rr[i].QType {
			t.Errorf("TC[%d] Failed to parse qtype : Wanted [%s] but found [%s]", i, rr[i].QType, l.QType)
		}
		if l.Query != rr[i].Query {
			t.Errorf("TC[%d] Failed to parse query : Wanted [%s] but found [%s]", i, rr[i].Query, l.Query)
		}
	}
}
