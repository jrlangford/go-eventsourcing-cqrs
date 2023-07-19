#!/bin/bash
UUID=$1
if [ -z "$UUID" ]
then
    echo "ERROR: UUID is empty"
    exit 1
fi
grpcurl \
    -import-path ../../../../schema/proto/ \
    -proto query.proto \
    -plaintext \
    -d @ \
    localhost:50052 \
    query.Query/Run <<EOM
    {
        "getInventoryItemDetails": {
            "uuid":"$UUID"
        }
    }
EOM
