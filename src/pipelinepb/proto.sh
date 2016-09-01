#!/bin/bash

if [ "$ROOT_DIR" = "" ];then
    ROOT_DIR="."
fi

VENDOR_PATH="$ROOT_DIR/vendor/src"
CWD=$ROOT_DIR/src/pipelinepb

export IMPORT_PATH=$VENDOR_PATH:$CWD
export GENERATOR="gofast_out"
export OUTPUT_DIR="$CWD"
export PROTO_FILES="$CWD/*.proto"
export OPTIONS_PROTO="${OPTIONS_API}"

protoc --proto_path=${IMPORT_PATH} \
    --${GENERATOR}=${OPTIONS_PROTO},plugins=grpc:${OUTPUT_DIR} \
    --grpc-gateway_out=${OUTPUT_DIR} \
    ${PROTO_FILES}
