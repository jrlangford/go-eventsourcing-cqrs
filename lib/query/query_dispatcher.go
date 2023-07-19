package query

import (
	"fmt"
	"reflect"
)

// QueryHandler wraps a Handle method for QueryMessages
type QueryHandler interface {
	Handle(QueryMessage) (interface{}, error)
}

// QueryDispatcher wraps the methods needed to register nd dispatch queries.
type QueryDispatcher interface {
	Dispatch(QueryMessage) error
	RegisterHandler(QueryHandler, ...interface{}) error
}

// An InMemoryQueryDispatcher is self-descriptive.
type InMemoryQueryDispatcher struct {
	handlers map[string]QueryHandler
}

// NewInMemoryQueryDispatcher is self-descriptive.
func NewInMemoryQueryDispatcher() *InMemoryQueryDispatcher {
	b := &InMemoryQueryDispatcher{
		handlers: make(map[string]QueryHandler),
	}
	return b
}

// Dispatch runs a Query using its corresponding handler, if present.
func (b *InMemoryQueryDispatcher) Dispatch(query QueryMessage) (interface{}, error) {
	if handler, ok := b.handlers[query.QueryType()]; ok {
		return handler.Handle(query)
	}
	return nil, fmt.Errorf("The query dispatcher does not have a handler for queries of type: %s", query.QueryType())
}

// RegisterHandler registers a query handler for one or more query types.
func (b *InMemoryQueryDispatcher) RegisterHandler(handler QueryHandler, queries ...interface{}) error {
	for _, query := range queries {
		typeName := reflect.TypeOf(query).Elem().Name()
		if _, ok := b.handlers[typeName]; ok {
			return fmt.Errorf("Duplicate query handler registration with query bus for query of type: %s", typeName)
		}
		b.handlers[typeName] = handler
	}
	return nil
}
