#!/bin/bash
set -e # script terminates when error occurs
BASE_PATH=$(
    cd "$(dirname "$0")" || exit
    pwd
)
CODE_PATH=${BASE_PATH}/..
cd "${CODE_PATH}/bin"
BIN_DIR=$(pwd)
cd "${CODE_PATH}"

while getopts "e:" arg; do
    # shellcheck disable=SC2220
    case ${arg} in
    e)
        build_env=${OPTARG}
        ;;
    esac
done

function go_install() {
    if [[ -z "${build_env}" || "${build_env}" == "local" ]]; then
        if [[ ! -x "$(command -v $1)" ]]; then
            echo "    [Install $1 to GOBIN. Or GOPATH if GOBIN is not set]"
            go install "$2"
        fi
    elif [[ "${build_env}" == "project" ]]; then
        if [[ ! -e "${BIN_DIR}/$1" ]]; then
            echo "    [Install $1 to ${BIN_DIR}/]"
            GOBIN="${BIN_DIR}" go install "$2"
        fi
    else
        echo "Unknown build env: ${build_env}"
        return 1
    fi
}

if [[ -n "${build_env}" && "${build_env}" != "local" && "${build_env}" != "project" ]]; then
    echo "Unknown build env: ${build_env}"
    exit 1
fi

go_install gocov github.com/axw/gocov/gocov@v1.0.0
go_install protoc-gen-validate github.com/envoyproxy/protoc-gen-validate@v0.6.1
go_install swagger github.com/go-swagger/go-swagger/cmd/swagger@v0.27.0
go_install mockgen github.com/golang/mock/mockgen@v1.6.0
go_install golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
go_install protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
go_install protoc-gen-openapiv2 github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.3.0
go_install protoc-gen-doc github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.4.1
go_install goimports golang.org/x/tools/cmd/goimports@v0.1.5
go_install protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
go_install protoc-gen-go github.com/golang/protobuf/protoc-gen-go@v1.4.0
