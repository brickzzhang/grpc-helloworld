#!/bin/bash
set -e # script terminates when error occurs

BASE_PATH=$(
    cd "$(dirname "$0")" || exit
    pwd
)
CODE_PATH=${BASE_PATH}/..
cd "${CODE_PATH}/bin"

BIN_DIR=$(pwd)
PROTOC_VERSION="3.15.6"
Kernel=$(uname -s)
Architecture=$(uname -m)

while getopts "e:" arg; do
    # shellcheck disable=SC2220
    case ${arg} in
    e)
        build_env=${OPTARG}
        ;;
    esac
done

if [[ -z "${build_env}" || "${build_env}" == "local" ]]; then
    INSTALL_LOC=""
elif [[ "${build_env}" == "project" ]]; then
    INSTALL_LOC="${BIN_DIR}/"
else
    echo "Unknown build env: ${build_env}"
    exit 1
fi

if [[ ! -e $(command -v "${INSTALL_LOC}protoc") ]]; then
    if [[ -z "${INSTALL_LOC}" ]]; then
        echo "    [Install protoc to /usr/local/bin]"
    else
        echo "    [Install protoc to ${INSTALL_LOC}]"
    fi
    if [[ ${Kernel} == "Darwin" ]]; then
        os_info="osx"
    elif [[ ${Kernel} == "Linux" ]]; then
        os_info="linux"
    else
        echo "OS not supported"
        exit 1
    fi
    PROTOC_ZIP="protoc-${PROTOC_VERSION}-${os_info}-${Architecture}.zip"
    protoc_url="https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"
    wget -q "${protoc_url}"

    unzip -q -o "${PROTOC_ZIP}" -d "/tmp" bin/protoc
    rm -f "${PROTOC_ZIP}"
    if [[ -z "${build_env}" || "${build_env}" == "local" ]]; then
        mv /tmp/bin/protoc /usr/local/bin
    elif [[ "${build_env}" == "project" ]]; then
        mv /tmp/bin/protoc "${INSTALL_LOC}"
    fi
    rm -rf /tmp/bin

    # check protoc
    if [[ ! -x $(command -v "${INSTALL_LOC}protoc") ]]; then
        echo "    [Error] protoc installation failed"
        exit 1
    fi
fi
