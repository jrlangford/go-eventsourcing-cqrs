// Package query implements features to build and dispatch query messages.
package query

import (
	"context"
	"reflect"
)

// QueryMessage wraps the query and its metadata.
type QueryMessage interface {
	Query() interface{}
	QueryType() string
	Context() context.Context
}

// QueryDescriptor is an implementation of the QueryMessage interface.
type QueryDescriptor struct {
	query interface{}
	ctx   context.Context
}

// NewQueryMessage returns a new QueryDescriptor.
func NewQueryMessage(ctx context.Context, query interface{}) *QueryDescriptor {
	return &QueryDescriptor{
		ctx:   ctx,
		query: query,
	}
}

// Query wraps specific query types.
func (q *QueryDescriptor) Query() interface{} {
	return q.query
}

// QueryType returns a string representation of the query type.
func (q *QueryDescriptor) QueryType() string {
	return reflect.TypeOf(q.query).Elem().Name()
}

// Context returns the context passed to the query upon creation.
func (q *QueryDescriptor) Context() context.Context {
	return q.ctx
}
