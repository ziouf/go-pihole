#! /usr/bin/env bash

# Install prerequisites
sudo apt update
sudo apt install -y build-essential nettle-dev libidn11-dev libdbus-1-dev libnetfilter-conntrack-dev
# Clone DNSMASQ sources
git clone http://thekelleys.org.uk/git/dnsmasq.git /usr/lib/dnsmasq
# Build DNSMASQ
cd /usr/lib/dnsmasq
sudo make all-i18n COPTS="-DHAVE_DNSSEC -DHAVE_DNSSEC_STATIC -DHAVE_DBUS -DHAVE_IDN -DHAVE_CONNTRACK --static"

sudo ln -s /usr/lib/dnsmasq/src/dnsmasq /usr/bin/go-pihole/dnsmasq


