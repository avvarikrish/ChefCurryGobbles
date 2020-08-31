#!/bin/bash

protoc proto/"$1"/"$1".proto --go_out=plugins=grpc:.
