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

// EventMessage is the interface that a command must implement.
type EventMessage interface {

	// AggregateID returns the ID of the Aggregate that the event relates to.
	AggregateID() string

	// Returns the actual event which is the payload of the event message.
	Event() interface{}

	// EventType returns a string descriptor of the command name.
	EventType() string

	// Version returns the version of the event.
	Version() *int

	// UserMetadata returns event metadata.
	UserMetadata() []byte

	// SetUserMetadata sets event metadata.
	SetUserMetadata([]byte)
}

// EventDescriptor is an implementation of the event message interface.
type EventDescriptor struct {
	aggregateID  string
	event        interface{}
	version      *int
	userMetadata []byte
}

// NewEventMessage returns a new event descriptor
func NewEventMessage(aggregateID string, event interface{}, version *int) *EventDescriptor {
	return &EventDescriptor{
		aggregateID:  aggregateID,
		event:        event,
		version:      version,
		userMetadata: []byte{},
	}
}

// EventType returns the name of the event type as a string.
func (c *EventDescriptor) EventType() string {
	return typeOf(c.event)
}

// AggregateID returns the ID of the Aggregate that the event relates to.
func (c *EventDescriptor) AggregateID() string {
	return c.aggregateID
}

// Event the event payload of the event message
func (c *EventDescriptor) Event() interface{} {
	return c.event
}

// Version returns the version of the event
func (c *EventDescriptor) Version() *int {
	return c.version
}

func (c *EventDescriptor) SetUserMetadata(metadata []byte) {
	c.userMetadata = metadata
}

func (c *EventDescriptor) UserMetadata() []byte {
	return c.userMetadata
}
