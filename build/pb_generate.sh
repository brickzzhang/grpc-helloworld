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

while getopts "e:p:" arg; do
    # shellcheck disable=SC2220
    case ${arg} in
    e)
        build_env=${OPTARG}
        ;;
    p)
        PROTOFILE_DIR=${OPTARG}
        ;;
    esac
done

# use default value
if [ -z "$PROTOFILE_DIR" ]; then
    PROTOFILE_DIR="api"
fi

echo "Executing protocol buffer code and swagger generation commands: [4 in total]"

# proto_generation generate proto code using proto file
function proto_generation() {
    GOOGLE_PB_PKG_PATH="${CODE_PATH}/api"

    if [[ -z "${build_env}" || "${build_env}" == "local" ]]; then
        echo "    [1/4] Generate service code from proto file at $1 ..."
        protoc -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --go_out=. --go-grpc_out=. $1
        echo "    [2/4] Generate reverse proxy from proto file at $1 ..."
        protoc -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --grpc-gateway_out=logtostderr=true:. $1
        echo "    [3/4] Generate validator"
        protoc -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --validate_out=lang=go:. $1
        echo "    [4/4] Compile swagger file at $1"
        protoc -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --openapiv2_out=logtostderr=true:. --openapiv2_opt=json_names_for_fields=false $1
    else
        echo "    [1/4] Generate service code from proto file at $1 ..."
        "protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-go="${BIN_DIR}/protoc-gen-go" --plugin=protoc-gen-go-grpc="${BIN_DIR}/protoc-gen-go-grpc" --go_out=. --go-grpc_out=. $1
        echo "    [2/4] Generate reverse proxy from proto file at $1 ..."
        "protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-grpc-gateway="${BIN_DIR}/protoc-gen-grpc-gateway" --grpc-gateway_out=logtostderr=true:. $1
        echo "    [3/4] Generate validator"
        "protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-validate="${BIN_DIR}/protoc-gen-validate" --validate_out=lang=go:. $1
        echo "    [4/4] Compile swagger file at $1"
        "protoc" -I. -I api -I "${GOOGLE_PB_PKG_PATH}/googleapis" -I api/third_party --plugin=protoc-gen-openapiv2="${BIN_DIR}/protoc-gen-openapiv2" --openapiv2_out=logtostderr=true:. --openapiv2_opt=json_names_for_fields=false $1
    fi
}

# read_proto_dir_and_generate traverse all proto files under specified dir
function read_proto_dir_and_generate() {
    for file in $(ls $1 | grep -v "^google$" | grep -v "^googleapis$" | grep -v "^third_party$"); do
        # file is a directory
        if [ -d $1"/"$file ]; then
            read_proto_dir_and_generate $1"/"$file
        # file is not a directory and not a proto file
        else
            if [ "$(ls -1 -F $1"/"$file | grep -c ".proto$")" -eq 0 ]; then
                echo "---------------not proto file: $1"/"$file---------------"
                continue
            fi

            proto_file=$(ls -1 -F $1"/"$file)
            proto_generation $proto_file
        fi
    done
}

# generate codes
if [[ -z "${build_env}" || "${build_env}" == "local" ]]; then
    echo "[Using tools in PATH of local machine]"

    # bash "${CODE_PATH}/build/protoc_install.sh" -e local
    bash "${CODE_PATH}/build/addons_install.sh" -e local

    read_proto_dir_and_generate $PROTOFILE_DIR
elif [[ "${build_env}" == "project" ]]; then
    echo "[Using tools in project bin: ${BIN_DIR}]"

    # bash "${CODE_PATH}/build/protoc_install.sh" -e project
    bash "${CODE_PATH}/build/addons_install.sh" -e project

    read_proto_dir_and_generate $PROTOFILE_DIR
else
    echo "Unknown build env: ${build_env}"
    exit 1
fi

# move_apigen_content_2_current_dir move all proto generated files under specified directory to apigen dir
function move_apigen_content_2_current_dir() {
    # get module name
    ARR=($(sed -n '/^module /p' go.mod))
    MODULE=${ARR[1]}
    echo "module is: $MODULE"

    # if go_package option is not defined in proto file, exit
    if [ ! -d $MODULE ]; then
        return 0
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
}

move_apigen_content_2_current_dir

function move_generated_swagger_json_2_swagger_dir() {
    if [ -d "${CODE_PATH}/api/swagger" ]; then
        rm -rf "${CODE_PATH}/api/swagger"
    fi
    mkdir -p "${CODE_PATH}/api/swagger"

    if [ "$PROTOFILE_DIR" = "api" ]; then
        find $PROTOFILE_DIR/* -name '*.swagger.json' | xargs -I '{}' cp {} "${CODE_PATH}/api/swagger"
        find $PROTOFILE_DIR -maxdepth 1 -name '*.swagger.json' | xargs rm -rf
    else
        # copy all swagger file to api/swagger
        if [ -n "$(ls $PROTOFILE_DIR/*)" ]; then
            cp -r $PROTOFILE_DIR/* ./api/swagger
        fi

        find "${CODE_PATH}/api/swagger" -type f -not -name '*.swagger.json' | xargs rm -rf
        find "$PROTOFILE_DIR" -name '*.swagger.json' | xargs rm -rf
    fi
}

move_generated_swagger_json_2_swagger_dir
