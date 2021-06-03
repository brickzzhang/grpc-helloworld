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

go_install gocov github.com/axw/gocov/gocov@latest
go_install protoc-gen-validate github.com/envoyproxy/protoc-gen-validate@latest
go_install swagger github.com/go-swagger/go-swagger/cmd/swagger@latest
go_install mockgen github.com/golang/mock/mockgen@latest
go_install golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go_install protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest
go_install protoc-gen-openapiv2 github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go_install protoc-gen-doc github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
go_install goimports golang.org/x/tools/cmd/goimports@latest
go_install protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go_install protoc-gen-go github.com/golang/protobuf/protoc-gen-go@v1.4.0
