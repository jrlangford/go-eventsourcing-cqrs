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

import (
	"fmt"
)

// StreamNamer is the interface that stream name delegates should implement.
type StreamNamer interface {
	GetStreamName(string) (string, error)
}

// AggregateStreamNamer stores delegates per aggregate type allowing fine grained
// control of stream names for event streams.
type AggregateStreamNamer struct {
	aggregateName string
}

// NewDelegateStreamNamer constructs a delegate stream namer
func NewAggregateStreamNamer(aggregateName string) *AggregateStreamNamer {
	return &AggregateStreamNamer{
		aggregateName: aggregateName,
	}
}

// GetStreamName gets the result of the stream name delgate registered for the aggregate type.
func (n *AggregateStreamNamer) GetStreamName(id string) (string, error) {
	if n.aggregateName == "" {
		return "", fmt.Errorf("aggregate name cannot be empty")
	}
	return fmt.Sprintf("%s_%s", n.aggregateName, id), nil
}
