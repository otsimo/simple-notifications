#!/bin/bash

ROOT_DIR=.
IMPORT_PATH=${ROOT_DIR}/vendor/src/:${ROOT_DIR}/notificationpb
OUTPUT=${ROOT_DIR}/notificationpb
PROTO_FILES=${ROOT_DIR}/notificationpb/service.proto

protoc \
    --proto_path=${IMPORT_PATH} \
    --gofast_out=plugins=grpc:${OUTPUT} \
    ${PROTO_FILES}

