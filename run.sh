#!  /usr/bin/env bash

DIR=$(cd $(dirname ${0}) && pwd)

# dep ensure
# dep ensure -update
# go get -u github.com/gorilla/handlers
# go get -u github.com/gorilla/mux
# go get -u github.com/hpcloud/tail
# go get -u github.com/spf13/viper
# go get -u github.com/jinzhu/gorm/...

# cd ${DIR}/ui

# npm run-script build

# cd ${DIR}

echo "========================="
echo "Starting go server"
go run *.go --dnsmasq-config-dir $(pwd)/tmp --dnsmasq-log-file $(pwd)/tmp/pihole.log