module github.com/jrlangford/go-eventsourcing-cqrs/services/command

go 1.20

require (
	github.com/EventStore/EventStore-Client-Go/v3 v3.0.0
	github.com/google/uuid v1.3.0
	github.com/jrlangford/go-eventsourcing-cqrs/lib v0.1.0
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/jrlangford/go-eventsourcing-cqrs/lib => ../../lib

require (
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20211116061358-0a5406a5449c // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
)
