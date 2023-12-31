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

// EventDispatcher is the inteface that an event bus must implement.
type EventDispatcher interface {
	Dispatch(EventMessage)
	AddHandler(EventHandler, ...interface{})
}

// InternalEventDispatcher provides a lightweight in process event bus
type InternalEventDispatcher struct {
	eventHandlers map[string]map[EventHandler]struct{}
}

// NewInternalEventDispatcher constructs a new InternalEventDispatcher
func NewInternalEventDispatcher() *InternalEventDispatcher {
	b := &InternalEventDispatcher{
		eventHandlers: make(map[string]map[EventHandler]struct{}),
	}
	return b
}

// Dispatch publishes events to all registered event handlers
func (b *InternalEventDispatcher) Dispatch(event EventMessage) {
	if handlers, ok := b.eventHandlers[event.EventType()]; ok {
		for handler := range handlers {
			handler.Handle(event)
		}
	}
}

// AddHandler registers an event handler for all of the events specified in the
// variadic events parameter.
func (b *InternalEventDispatcher) AddHandler(handler EventHandler, events ...interface{}) {

	for _, event := range events {
		typeName := typeOf(event)

		// There can be multiple handlers for any event.
		// Here we check that a map is initialized to hold these handlers
		// for a given type. If not we create one.
		if _, ok := b.eventHandlers[typeName]; !ok {
			b.eventHandlers[typeName] = make(map[EventHandler]struct{})
		}

		// Add this handler to the collection of handlers for the type.
		b.eventHandlers[typeName][handler] = struct{}{}
	}
}
