#!  /usr/bin/env bash

DIR=$(cd $(dirname ${0}) && pwd)
GODIR=${DIR}/backend
UIDIR=${DIR}/frontend

GOPACKAGES=cm-cloud.fr/go-pihole/backend/...

fn_dep() {
    echo "========================="
    echo "Update depecendies"
    cd ${GODIR}

    dep ensure -v
    # dep ensure -update -v

    cd ${DIR}
}

fn_build() {
    echo "========================="
    echo "Build Frontend"
    cd ${UIDIR}

    npm run-script build
    
    cd ${DIR}

    echo "========================="
    echo "Build Backend"
    cd ${GODIR}

    export CGO_ENABLED=0
    for GOARM in 5 6 7
    do 
        export GOOS=linux
        export GOARCH=arm
        export GOARM
        echo "Building $GOOS-$GOARCH-$GOARM"
        go build -o ${GODIR}/bin/${GOOS}/${GOARCH}-${GOARM}/go-pihole -a -ldflags '-extldflags "-static"'
    done

    for GOOS in darwin linux windows; do
        for GOARCH in 386 amd64; do
            export GOOS 
            export GOARCH
            echo "Building $GOOS-$GOARCH"
            go build -o ${GODIR}/bin/${GOOS}/${GOARCH}/go-pihole -a -ldflags '-extldflags "-static"'
        done
    done
    
    cd ${DIR}
}

fn_package() {
    echo "========================="
    echo "Package"
    package_file=go-pihole.tar.gz
 
    [ -f ${DIR}/${package_file} ] && rm ${DIR}/${package_file}

    tar czf ${DIR}/${package_file} backend/bin frontend/dist
}

fn_test() {
    echo "========================="
    echo "Run go tests"
    go test ${GOPACKAGES}
}

fn_run() {
    echo "========================="
    echo "Starting go server"
    go run ${GODIR}/*.go                        \
    --db.cleaning.enable                        \
    --dnsmasq.config.dir $(pwd)/tmp             \
    --dnsmasq.log.file $(pwd)/tmp/pihole.log    \
    --log.level VERBOSE                         \
    --log.path $(pwd)/tmp                       \
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
        fn_dep
        fn_test
        fn_build
        fn_package
    ;;
    run)
        fn_test
        fn_run
    ;;
    *)
        fn_usage
    ;;
esac
