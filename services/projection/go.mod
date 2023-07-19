module github.com/jrlangford/go-eventsourcing-cqrs/services/projection

go 1.20

require (
	github.com/EventStore/EventStore-Client-Go/v3 v3.0.0
	github.com/google/uuid v1.1.2
	github.com/jrlangford/go-eventsourcing-cqrs/lib v0.1.0
	github.com/redis/go-redis/v9 v9.0.5
)

replace github.com/jrlangford/go-eventsourcing-cqrs/lib => ../../lib

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20211116061358-0a5406a5449c // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
	google.golang.org/grpc v1.45.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
