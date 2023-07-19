// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE_2016 file.

// Modifications notice:
// Base version: https://github.com/jetbasrawi/go.cqrs/commit/e4d812d57f090ecede016aa36d70c73626a8eb17
// Copyright 2023 Jonathan Langford.
//
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file at the root of this project.

package eventsourcing

// AggregateFactory returns aggregate instances of a specified type with the
// AggregateID set to the uuid provided.
//
// An aggregate factory is typically a dependency of the repository that will
// delegate instantiation of aggregate instances to the Aggregate factory.
type AggregateFactory[D any, T interface {
	*D
	AggregateRoot
}] interface {
	GetAggregate(string) T
}

// DelegateAggregateFactory is an implementation of the AggregateFactory interface
// that supports registration of delegate functions to perform aggregate instantiation.
type DelegateAggregateFactory[D any, T interface {
	*D
	AggregateRoot
}] struct {
	constructor func(string) T
}

// NewDelegateAggregateFactory contructs a new DelegateAggregateFactory
func NewDelegateAggregateFactory[D any, T interface {
	*D
	AggregateRoot
}](constructor func(string) T) *DelegateAggregateFactory[D, T] {
	return &DelegateAggregateFactory[D, T]{
		constructor: constructor,
	}
}

// GetAggregate calls the delegate for the type specified and returns the result.
func (f *DelegateAggregateFactory[D, T]) GetAggregate(id string) T {
	return f.constructor(id)
}
