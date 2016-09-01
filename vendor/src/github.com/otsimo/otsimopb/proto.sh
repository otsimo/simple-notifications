#!/bin/bash

export IMPORT_PATH=$GOPATH/src:.
export GENERATOR="gogofaster_out"
export OUTPUT_DIR="."
export PROTO_FILES="./*.proto"

protoc --proto_path=$IMPORT_PATH --${GENERATOR}=plugins=grpc:${OUTPUT_DIR} $PROTO_FILES

