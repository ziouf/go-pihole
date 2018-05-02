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
	const logLine = `Apr  9 00:50:24 dnsmasq[755]: 78348 192.168.1.254/52101 query[A] a.root-servers.net from 192.168.1.254`

	l := NewLog().ParseLine(logLine)

	if date := time.Date(int(0), time.April, int(9), int(0), int(50), int(24), int(0), time.UTC); l.Date != date {
		t.Errorf("Failed to parse date : Wanted [%s] but found [%s]", date, l.Date)
		t.Fail()
	}
	if proc := "dnsmasq"; l.Process != proc {
		t.Errorf("Failed to parse pocess : Wanted [%s] but found [%s]", proc, l.Process)
		t.Fail()
	}
	if pid := int(755); l.PID != pid {
		t.Errorf("Failed to parse pid : Wanted [%d] but found [%d]", pid, l.PID)
		t.Fail()
	}
	if qid := int(78348); l.QID != qid {
		t.Errorf("Failed to parse qid : Wanted [%d] but found [%d]", qid, l.QID)
		t.Fail()
	}
	if ip := "192.168.1.254"; l.SourceIP != ip {
		t.Errorf("Failed to parse source ip : Wanted [%s] but found [%s]", ip, l.SourceIP)
		t.Fail()
	}
	if qtype := "query[A]"; l.QType != qtype {
		t.Errorf("Failed to parse qtype : Wanted [%s] but found [%s]", qtype, l.QType)
		t.Fail()
	}
	if query := "a.root-servers.net from 192.168.1.254"; l.Query != query {
		t.Errorf("Failed to parse query : Wanted [%s] but found [%s]", query, l.Query)
		t.Fail()
	}
}

func testParseDnsmasqDhcpLogLine(t *testing.T) {
	const logLine = `Apr  9 00:52:33 dnsmasq-dhcp[755]: L'option dhcp-option redondante 3 sera ignorée`

	l := NewLog().ParseLine(logLine)

	if date := time.Date(int(0), time.April, int(9), int(0), int(52), int(33), int(0), time.UTC); l.Date != date {
		t.Errorf("Failed to parse date : Wanted [%s] but found [%s]", date, l.Date)
		t.Fail()
	}
	if proc := "dnsmasq-dhcp"; l.Process != proc {
		t.Errorf("Failed to parse date : Wanted [%s] but found [%s]", proc, l.Process)
		t.Fail()
	}
	if pid := int(755); l.PID != pid {
		t.Errorf("Failed to parse pid : Wanted [%d] but found [%d]", pid, l.PID)
		t.Fail()
	}
	if query := "L'option dhcp-option redondante 3 sera ignorée"; l.Query != query {
		t.Errorf("Failed to parse pid : Wanted [%s] but found [%s]", query, l.Query)
		t.Fail()
	}


}