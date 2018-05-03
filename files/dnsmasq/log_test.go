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
		`Apr  9 20:13:29 dnsmasq[755]: 172814 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35/50227 reply image01.bonprix.fr is <CNAME>`,
		`Apr  9 20:13:29 dnsmasq[755]: 172814 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35/50227 validation result is INSECURE`,
		`Apr  9 00:52:36 dnsmasq[755]: 78395 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35/32969 query[AAAA] ads.nexage.com from 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`,
		`Apr  9 00:52:36 dnsmasq[755]: 78395 2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35/32969 /etc/pihole/gravity.list ads.nexage.com is 2a01:cb15:8046:3200:5ddd:efb1:d506:7ed6`,
	}
	rr := []Log{
		{
			Date: time.Date(int(0), time.April, int(9), int(0), int(50), int(24), int(0), time.UTC), Process: `dnsmasq`, PID: 755,
			DnsQID : 78348, DnsFrom: `192.168.1.254`, DnsType: `query[A]`, DnsFQDN: `a.root-servers.net`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(38), int(0), time.UTC), Process: `dnsmasq`, PID: 755,
			DnsQID : 78399, DnsFrom: `192.168.1.200`, DnsType: `cached`, DnsFQDN: `cloudsync-tw.synology.com`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(13), int(29), int(0), time.UTC), Process: `dnsmasq`, PID: 755,
			DnsQID : 172814, DnsFrom: `2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`, DnsType: `reply`, DnsFQDN: `image01.bonprix.fr`, DnsIP: `<CNAME>`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(13), int(29), int(0), time.UTC), Process: `dnsmasq`, PID: 755,
			DnsQID: 172814, DnsFrom: `2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`, DnsType: `validation`, DnsSecured: false,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(36), int(0), time.UTC), Process: `dnsmasq`, PID: 755,
			DnsQID: 78395, DnsFrom: `2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`, DnsType: `query[AAAA]`, DnsFQDN: `ads.nexage.com`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(0), int(52), int(36), int(0), time.UTC), Process: `dnsmasq`, PID: 755,
			DnsQID: 78395, DnsFrom: `2a01:cb15:8046:3200:8c3a:f31:9ac7:9c35`, DnsType: `/etc/pihole/gravity.list`, DnsFQDN: `ads.nexage.com`,
		},
	}

	for i, tc := range tcs {
		l := NewLog().ParseLine(tc)

		if l.Date != rr[i].Date {
			t.Errorf("TC[%d] Failed to parse Date : Wanted [%s] but found [%s]", i, rr[i].Date, l.Date)
		}
		if l.Process != rr[i].Process {
			t.Errorf("TC[%d] Failed to parse Process : Wanted [%s] but found [%s]", i, rr[i].Process, l.Process)
		}
		if l.PID != rr[i].PID {
			t.Errorf("TC[%d] Failed to parse PID : Wanted [%d] but found [%d]", i, rr[i].PID, l.PID)
		}
		if l.DnsQID != rr[i].DnsQID {
			t.Errorf("TC[%d] Failed to parse QID : Wanted [%d] but found [%d]", i, rr[i].DnsQID, l.DnsQID)
		}
		if l.DnsFrom != rr[i].DnsFrom {
			t.Errorf("TC[%d] Failed to parse From : Wanted [%s] but found [%s]", i, rr[i].DnsFrom, l.DnsFrom)
		}
		if l.DnsType != rr[i].DnsType {
			t.Errorf("TC[%d] Failed to parse Type : Wanted [%s] but found [%s]", i, rr[i].DnsType, l.DnsType)
		}
		if l.DnsFQDN != rr[i].DnsFQDN {
			t.Errorf("TC[%d] Failed to parse FQDN : Wanted [%s] but found [%s]", i, rr[i].DnsFQDN, l.DnsFQDN)
		}
		if l.DnsSecured != rr[i].DnsSecured {
			t.Errorf("TC[%d] Failed to parse Secured : Wanted [%t] but found [%t]", i, rr[i].DnsSecured, l.DnsSecured)
		}
	}
}

func testParseDnsmasqDhcpLogLine(t *testing.T) {
	tcs := []string{
		`Apr  9 20:14:13 dnsmasq-dhcp[755]: DHCPDISCOVER(eth0) a0:4c:5b:b3:72:50`,
		`Apr  9 20:14:13 dnsmasq-dhcp[755]: DHCPOFFER(eth0) 192.168.1.94 a0:4c:5b:b3:72:50`,
		`Apr  9 20:14:13 dnsmasq-dhcp[755]: DHCPREQUEST(eth0) 192.168.1.94 a0:4c:5b:b3:72:50`,
		`Apr  9 20:14:14 dnsmasq-dhcp[755]: DHCPACK(eth0) 192.168.1.94 a0:4c:5b:b3:72:50 android-496d52bb579c54f1`,
		`Apr  9 20:14:15 dnsmasq-dhcp[755]: RTR-SOLICIT(eth0) a0:4c:5b:b3:72:50`,
		`Apr  9 20:14:15 dnsmasq-dhcp[755]: RTR-ADVERT(eth0) 2a01:cb15:8046:3200::`,
		`Apr  9 20:14:15 dnsmasq-dhcp[755]: L'option dhcp-option redondante 23 sera ignorée`,
		`Apr  9 20:14:16 dnsmasq-dhcp[755]: SLAAC-CONFIRM(eth0) 2a01:cb15:8046:3200:a24c:5bff:feb3:7250 android-496d52bb579c54f1`,
	}
	rr := []Log{
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(13), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `DHCPDISCOVER(eth0)`, DhcpMAC: `a0:4c:5b:b3:72:50`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(13), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `DHCPOFFER(eth0)`, DhcpIP: `192.168.1.94`, DhcpMAC: `a0:4c:5b:b3:72:50`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(13), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `DHCPREQUEST(eth0)`, DhcpIP: `192.168.1.94`, DhcpMAC: `a0:4c:5b:b3:72:50`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(14), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `DHCPACK(eth0)`, DhcpIP: `192.168.1.94`, DhcpMAC: `a0:4c:5b:b3:72:50`, DhcpName: `android-496d52bb579c54f1`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(15), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `RTR-SOLICIT(eth0)`, DhcpMAC: `a0:4c:5b:b3:72:50`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(15), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `RTR-ADVERT(eth0)`, DhcpIP: `2a01:cb15:8046:3200::`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(15), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpMsg: `L'option dhcp-option redondante 23 sera ignorée`,
		},
		{
			Date: time.Date(int(0), time.April, int(9), int(20), int(14), int(16), int(0), time.UTC), Process: `dnsmasq-dhcp`, PID: 755,
			DhcpType: `SLAAC-CONFIRM(eth0)`, DhcpIP: `2a01:cb15:8046:3200:a24c:5bff:feb3:7250`, DhcpName: `android-496d52bb579c54f1`,
		},
	}

	for i, tc := range tcs {
		l := NewLog().ParseLine(tc)

		if l.Date != rr[i].Date {
			t.Errorf("TC[%d] Failed to parse Date : Wanted [%s] but found [%s]", i, rr[i].Date, l.Date)
		}
		if l.Process != rr[i].Process {
			t.Errorf("TC[%d] Failed to parse Pocess : Wanted [%s] but found [%s]", i, rr[i].Process, l.Process)
		}
		if l.PID != rr[i].PID {
			t.Errorf("TC[%d] Failed to parse PID : Wanted [%d] but found [%d]", i, rr[i].PID, l.PID)
		}
		if l.DhcpType != rr[i].DhcpType {
			t.Errorf("TC[%d] Failed to parse Type : Wanted [%s] but found [%s]", i, rr[i].DhcpType, l.DhcpType)
		}
		if l.DhcpMAC != rr[i].DhcpMAC {
			t.Errorf("TC[%d] Failed to parse MAC : Wanted [%s] but found [%s]", i, rr[i].DhcpMAC, l.DhcpMAC)
		}
		if l.DhcpIP != rr[i].DhcpIP {
			t.Errorf("TC[%d] Failed to parse IP : Wanted [%s] but found [%s]", i, rr[i].DhcpIP, l.DhcpIP)
		}
		if l.DhcpName != rr[i].DhcpName {
			t.Errorf("TC[%d] Failed to parse Name : Wanted [%s] but found [%s]", i, rr[i].DhcpName, l.DhcpName)
		}
		if l.DhcpMsg != rr[i].DhcpMsg {
			t.Errorf("TC[%d] Failed to parse Msg : Wanted [%s] but found [%s]", i, rr[i].DhcpMsg, l.DhcpMsg)
		}
	}
}
