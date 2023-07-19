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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
)

// // DomainRepository is the interface that all domain repositories should implement.
type AggregateRepoBase[D any, T interface {
	*D
	AggregateRoot
}] struct {
	eventStore         *esdb.Client
	streamNameDelegate StreamNamer
	aggregateFactory   AggregateFactory[D, T]
	eventFactory       EventFactory
}

// NewCommonDomainRepository constructs a new CommonDomainRepository
func NewAggregateRepoBase[D any, T interface {
	*D
	AggregateRoot
}](eventStore *esdb.Client) (*AggregateRepoBase[D, T], error) {
	if eventStore == nil {
		return nil, fmt.Errorf("Nil Eventstore injected into repository.")
	}

	d := &AggregateRepoBase[D, T]{
		eventStore: eventStore,
	}
	return d, nil
}

// SetAggregateFactory sets the aggregate factory that should be used to
// instantate aggregate instances
//
// Only one AggregateFactory can be registered at any one time.
// Any registration will overwrite the provious registration.
func (r *AggregateRepoBase[D, T]) SetAggregateFactory(factory AggregateFactory[D, T]) {
	r.aggregateFactory = factory
}

// SetEventFactory sets the event factory that should be used to instantiate event
// instances.
//
// Only one event factory can be set at a time. Any subsequent registration will
// overwrite the previous factory.
func (r *AggregateRepoBase[D, T]) SetEventFactory(factory EventFactory) {
	r.eventFactory = factory
}

// SetStreamNameDelegate sets the stream name delegate
func (r *AggregateRepoBase[D, T]) SetStreamNameDelegate(delegate StreamNamer) {
	r.streamNameDelegate = delegate
}

// Load will load all events from a stream and apply those events to an aggregate
// of the type specified.
//
// The aggregate type and id will be passed to the configured StreamNamer to
// get the stream name.
func (r *AggregateRepoBase[D, T]) Load(aggregateID string) (T, error) {

	if r.aggregateFactory == nil {
		return nil, fmt.Errorf("The repository has no Aggregate Factory.")
	}

	if r.streamNameDelegate == nil {
		return nil, fmt.Errorf("The repository has no stream name delegate.")
	}

	if r.eventFactory == nil {
		return nil, fmt.Errorf("The repository has no Event Factory.")
	}

	aggregate := r.aggregateFactory.GetAggregate(aggregateID)
	if aggregate == nil {
		return nil, fmt.Errorf("The returned aggregate is nil.")
	}

	streamName, err := r.streamNameDelegate.GetStreamName(aggregateID)
	if err != nil {
		return nil, err
	}

	stream, err := r.eventStore.ReadStream(context.Background(), streamName, esdb.ReadStreamOptions{}, 10)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	for {

		streamedEventMessage, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		streamedEvent := streamedEventMessage.Event

		domainEvent, err := r.eventFactory.GetEvent(streamedEvent.EventType, streamedEvent.Data)
		if err != nil {
			return nil, err
		}

		em := NewEventMessage(
			aggregateID,
			domainEvent,
			Int(int(streamedEvent.EventNumber)),
		)
		em.SetUserMetadata(streamedEvent.UserMetadata)

		aggregate.Apply(em, false)
		aggregate.IncrementVersion()
	}

	return aggregate, nil

}

// Save persists an aggregate
func (r *AggregateRepoBase[D, T]) Save(aggregate AggregateRoot, expectedVersion *int) error {

	if r.streamNameDelegate == nil {
		return fmt.Errorf("The repository has no stream name delagate.")
	}

	resultEvents := aggregate.GetChanges()

	streamName, err := r.streamNameDelegate.GetStreamName(aggregate.AggregateID())
	if err != nil {
		return err
	}

	if len(resultEvents) > 0 {

		evs := make([]esdb.EventData, len(resultEvents))

		for i, v := range resultEvents {
			evJson, err := json.Marshal(v.Event())
			if err != nil {
				log.Fatal(err)
			}

			evs[i].ContentType = esdb.ContentTypeJson
			evs[i].EventType = v.EventType()
			evs[i].Data = evJson
			evs[i].Metadata = v.UserMetadata()
		}

		aopts := esdb.AppendToStreamOptions{
			ExpectedRevision: esdb.Revision(uint64(aggregate.OriginalVersion())),
		}
		_, err = r.eventStore.AppendToStream(context.Background(), streamName, aopts, evs...)
		if err != nil {
			log.Fatal(err)
		}
	}

	aggregate.ClearChanges()

	return nil
}
