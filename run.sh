#!  /usr/bin/env bash

DIR=$(cd $(dirname ${0}) && pwd)

fn_dep() {
    echo "========================="
    echo "Update depecendies"

    # dep ensure
    # dep ensure -update

    go get -u github.com/gorilla/handlers
    go get -u github.com/gorilla/mux
    go get -u github.com/hpcloud/tail
    go get -u github.com/spf13/viper
    go get -u github.com/jinzhu/gorm/...
}

fn_build() {
    echo "========================="
    echo "Build UI"
    cd ${DIR}/ui

    npm run-script build
    
    cd ${DIR}

    echo "========================="
    echo "Build Backend"

    go build
}

fn_run() {
    echo "========================="
    echo "Starting go server"
    go run *.go --dnsmasq-config-dir $(pwd)/tmp --dnsmasq-log-file $(pwd)/tmp/pihole.log
}

fn_install() {
    echo "========================="
    echo "Installing application on system"



}

fn_usage() {
    echo "$0 : ..."
}

case $1 in
    dep)
        fn_dep
    ;;
    build)
        fn_build
    ;;
    run)
        fn_run
    ;;
    *)
        fn_usage
    ;;
esac


# cd ${DIR}/ui

# npm run-script build

# cd ${DIR}
