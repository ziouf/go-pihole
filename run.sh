#!  /usr/bin/env bash

DIR=$(cd $(dirname ${0}) && pwd)

# dep ensure
# dep ensure -update

# cd ${DIR}/ui

# npm run-script build

# cd ${DIR}

echo "========================="
echo "Starting go server"
go run *.go --dnsmasq-config-dir $(pwd)/tmp