PROTOC= protoc \
	--proto_path=.schema/proto/ \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=$$GOOPTS \
	--go-grpc_opt=$$GOOPTS

COMMAND_PB_PATH=services/command/internal/primadapters/pb
COMMAND_PROTO=$(COMMAND_PB_PATH)/command.pb.go $(COMMAND_PB_PATH)/command_grpc.pb.go

QUERY_PB_PATH=services/query/internal/primadapters/pb
QUERY_PROTO=$(QUERY_PB_PATH)/query.pb.go $(QUERY_PB_PATH)/query_grpc.pb.go

UI_COMMAND_PB_PATH=services/ui/internal/server/pbcommand
UI_COMMAND_PROTO=$(UI_COMMAND_PB_PATH)/command.pb.go $(UI_COMMAND_PB_PATH)/command_grpc.pb.go

UI_QUERY_PB_PATH=services/ui/internal/server/pbquery
UI_QUERY_PROTO=$(UI_QUERY_PB_PATH)/query.pb.go $(UI_QUERY_PB_PATH)/query_grpc.pb.go

.PHONY: proto
proto: $(COMMAND_PROTO) $(UI_COMMAND_PROTO) $(QUERY_PROTO) $(UI_QUERY_PROTO)

$(COMMAND_PROTO) &: schema/proto/command.proto
	GOOPTS="Mcommand.proto=$(COMMAND_PB_PATH)" && \
	$(PROTOC) \
	command.proto

$(QUERY_PROTO) &: schema/proto/query.proto
	GOOPTS="Mquery.proto=$(QUERY_PB_PATH)" && \
	$(PROTOC) \
	query.proto

$(UI_COMMAND_PROTO) &: schema/proto/command.proto
	GOOPTS="Mcommand.proto=$(UI_COMMAND_PB_PATH)" && \
	$(PROTOC) \
	command.proto

$(UI_QUERY_PROTO) &: schema/proto/query.proto
	GOOPTS="Mquery.proto=$(UI_QUERY_PB_PATH)" && \
	$(PROTOC) \
	query.proto

.PHONY: docker-up
docker-up:
	docker-compose -p go-eventsourcing up -d

.PHONY: docker-down
docker-down:
	docker-compose -p go-eventsourcing down

.PHONY: redis-flushdb
redis-flushdb:
	docker exec go-eventsourcing-redis-1 redis-cli flushdb

.PHONY: eventstore-flushdb
eventstore-flushdb:
	docker restart go-eventsourcing-eventstore-1

.PHONY: flush
flush: redis-flushdb eventstore-flushdb
