#!/bin/sh

protoc -I=sio/pb --go_out=sio/pb --go_opt=paths=source_relative sio/pb/*.proto

