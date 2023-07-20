# Go-eventsourcing-cqrs

A reference project that implements event sourcing and CQRS, based on [Jet
Basrawi's implementation](https://github.com/jetbasrawi/go.cqrs).

The main goal of this implementation is to clearly show responsibility
segregation in an event sourced system at the architectural level by running
each responsibility in an independent microservice.

## Requirements

### Execution

- docker = 24.0
- docker buildx = 0.10
- docker-compose = 2.18

### Development

- make = 4.4
- golang = 1.20
- protoc = 23.3
- protoc-gen-go = 1.3
- protoc-gen-go-grpc = 1.3

## Quick start

Build docker images.

```bash
docker buildx bake
```

Run the system with docker-compose.

```bash
docker-compose up
```

Open `localhost:8088` in your browser and interact with the UI.

> Important: The UI is not reactive, you may need to refresh your browser after making an update.

## Services

### Command

Processes incoming command messages and executes them, writing the resulting events to EventStoreDB.

Command messages are received via RPC.

### Projection

Processes incoming events and writes the corresponding projections to Redis.

Incoming events are received through a persistent subscription to EventStoreDB.

### Query

Processes incoming query messages and runs them, returning the corresponding data from Redis.

Query messages are received via RPC.

### UI

Serves rendered html files that let the user interact with the system via forms.

Forms trigger RPC commands when submitted.
Html files are rendered with information fetched through RPC queries.
