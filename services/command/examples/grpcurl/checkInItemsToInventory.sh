#!/bin/bash
UUID=$1
if [ -z "$UUID" ]
then
    echo "ERROR: uuid is empty"
    exit 1
fi
COUNT=$2
if [ -z "$COUNT" ]
then
    echo "ERROR: count is empty"
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
    "checkInItemsToInventory": {
        "uuid":"$UUID",
        "count":"$COUNT"
    }
}
EOM
