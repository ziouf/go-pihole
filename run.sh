#!  /usr/bin/env bash

DIR=$(cd $(dirname ${0}) && pwd)

fn_dep() {
    echo "========================="
    echo "Update depecendies"

    # dep ensure
    dep ensure -update -v

    # go get -u github.com/gorilla/handlers
    # go get -u github.com/gorilla/mux
    # go get -u github.com/hpcloud/tail
    # go get -u github.com/spf13/viper
    # go get -u github.com/boltdb/bolt
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

fn_test() {
    echo "========================="
    echo "Run go tests"
    go test cm-cloud.fr/go-pihole/...
}

fn_run() {
    echo "========================="
    echo "Starting go server"
    go run ${DIR}/*.go                          \
    --db.cleaning.enable                        \
    --dnsmasq.config.dir $(pwd)/tmp             \
    --dnsmasq.log.file $(pwd)/tmp/pihole.log    \
    --log.level VERBOSE                         \
    --log.file $(pwd)/go-pihole.log             \
    # --dnsmasq.bin $(pwd)/../dnsmasq/src/dnsmasq \
    # --db.file $(pwd)/go-pihole.db           \
    # --dnsmasq.embeded
}

fn_install() {
    echo "========================="
    echo "Installing application on system"

    # TODO : Add install scripts 

}

fn_usage() {
    echo "$0 : ..."
}

case $1 in
    dep)
        fn_dep
    ;;
    test)
        fn_test
    ;;
    build)
        fn_test
        fn_build
    ;;
    run)
        fn_test
        fn_run
    ;;
    *)
        fn_usage
    ;;
esac
