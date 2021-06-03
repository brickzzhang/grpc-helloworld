#!/bin/bash
set -e
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

echo "Executing protocol buffer code and swagger generation commands: [3 in total]"

if [[ "${build_env}" == "project" ]]; then
    echo "[Using tools in project bin: ${BIN_DIR}]"

    bash "${CODE_PATH}/build/protoc_install.sh" -e project
    bash "${CODE_PATH}/build/addons_install.sh" -e project

    GOOGLE_PB_PKG_PATH="${CODE_PATH}/api"

    echo "    [1/3] Generate service code from proto file at api/*.proto ..."
    "${BIN_DIR}/protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-go="${BIN_DIR}/protoc-gen-go" --plugin=protoc-gen-go-grpc="${BIN_DIR}/protoc-gen-go-grpc" --go_out=. --go-grpc_out=. api/*.proto
    echo "    [2/3] Generate reverse proxy from proto file at api/*.proto ..."
    "${BIN_DIR}/protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-grpc-gateway="${BIN_DIR}/protoc-gen-grpc-gateway" --grpc-gateway_out=logtostderr=true:. api/*.proto
    echo "    [3/3] Generate validator"
    "${BIN_DIR}/protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-validate="${BIN_DIR}/protoc-gen-validate" --validate_out=lang=go:. api/*.proto
else
    echo "Unknown build env: ${build_env}"
    exit 1
fi

# get module name
ARR=($(sed -n '/^module /p' go.mod))
MODULE=${ARR[1]}
echo "module is: $MODULE"

# if go_package option is not defined in proto file, exit
if [ ! -d $MODULE ]; then
    exit 0
fi

# clean target dir
for dir in $(ls $MODULE); do
    if [ -d ./$dir ]; then
        rm -rf ./$dir
    else
        rm -f ./$dir
    fi
done

# move dir
for dir in $(ls $MODULE); do
    mv -f $MODULE/$dir ./
done

#  remove code dir
OLD_IFS="$IFS"
IFS="/"
ARR_MODULE=($MODULE)
DIR=${ARR_MODULE[0]}
IFS=$OLD_IFS
rm -rf $DIR
