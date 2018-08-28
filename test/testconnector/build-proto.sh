#!/usr/bin/env bash
protoc --proto_path=qlik_connect --go_out=plugins=grpc:qlik_connect    grpc_server.proto