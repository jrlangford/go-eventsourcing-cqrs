#!/bin/bash
NAME=$1
if [ -z "$NAME" ]
then
    echo "ERROR: name is empty"
    exit 1
fi
grpcurl \
    -import-path ../../../../schema/proto/ \
    -proto command.proto \
    -plaintext \
    -d @ \
    localhost:50051 \
    command.Command/Execute <<EOM
{
    "createInventoryItem": {
        "name":"$NAME"
    }
}
EOM
