#!/bin/bash
grpcurl \
    -import-path ../../../../schema/proto/ \
    -proto query.proto \
    -plaintext \
    -d @ \
    localhost:50052 \
    query.Query/Run <<EOM
    {
        "getInventoryItems":{
        }
    }
EOM
