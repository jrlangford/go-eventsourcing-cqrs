package primadapters

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/google/uuid"
	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/lib/jsontools"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/primports"
)

// An EventConsumer reads events from the event store and processes them.
type EventConsumer struct {
	esdbClient   *esdb.Client
	projector    primports.EventProjector
	subscription *esdb.PersistentSubscription
	eventFactory *es.DelegateEventFactory
}

// NewEventConsumer is self-describing.
func NewEventConsumer(esdbClient *esdb.Client, projector primports.EventProjector) *EventConsumer {

	eventFactory := es.NewDelegateEventFactory()
	eventFactory.RegisterDelegate(&domain.InventoryItemCreated{}, jsontools.JsonToTypeInterface[domain.InventoryItemCreated])
	eventFactory.RegisterDelegate(&domain.InventoryItemRenamed{}, jsontools.JsonToTypeInterface[domain.InventoryItemRenamed])
	eventFactory.RegisterDelegate(&domain.InventoryItemDeactivated{}, jsontools.JsonToTypeInterface[domain.InventoryItemDeactivated])
	eventFactory.RegisterDelegate(&domain.ItemsRemovedFromInventory{}, jsontools.JsonToTypeInterface[domain.ItemsRemovedFromInventory])
	eventFactory.RegisterDelegate(&domain.ItemsCheckedIntoInventory{}, jsontools.JsonToTypeInterface[domain.ItemsCheckedIntoInventory])

	return &EventConsumer{
		esdbClient:   esdbClient,
		projector:    projector,
		eventFactory: eventFactory,
	}
}

// InitializeSubscription creates a persistent subscription to the event store.
func (ec *EventConsumer) InitializeSubscription(groupName string) error {
	// TODO: handle error after fixing esdb client library. It currently returns an empty error, making it impossible to identify if the error indicates the group already exists.
	_ = ec.esdbClient.CreatePersistentSubscriptionToAll(
		context.Background(),
		groupName,
		esdb.PersistentAllSubscriptionOptions{},
	)

	sub, err := ec.esdbClient.SubscribeToPersistentSubscriptionToAll(
		context.Background(),
		groupName,
		esdb.SubscribeToPersistentSubscriptionOptions{},
	)
	if err == nil {
		ec.subscription = sub
	}

	return err
}

// Run consumes and processes events continuously.
func (ec *EventConsumer) Run() {
	for {
		activeSub := ec.subscription.Recv()

		if activeSub.EventAppeared != nil {
			recordedEvent := activeSub.EventAppeared.Event.Event
			fmt.Printf("Event:\n%+v\n------------\n", recordedEvent)

			if recordedEvent.EventType == "$metadata" {
				fmt.Printf("Metadata event: %+v\n", recordedEvent.Data)
				ec.subscription.Ack(activeSub.EventAppeared.Event)
			} else if recordedEvent.EventType == "PersistentConfig1" {
				fmt.Printf("Metadata event: %+v\n", recordedEvent.Data)
				ec.subscription.Ack(activeSub.EventAppeared.Event)
			} else {

				domainEvent, err := ec.eventFactory.GetEvent(recordedEvent.EventType, recordedEvent.Data)
				if err != nil {
					log.Fatal(err)
				}

				streamID := recordedEvent.StreamID
				aggregateID, err := getAggregateIDFromStreamID(streamID)
				if err != nil {
					log.Fatal(err)
				}

				err = ec.projector.ProjectEvent(
					aggregateID,
					domainEvent,
					es.Int(int(recordedEvent.EventNumber)),
					recordedEvent.UserMetadata,
				)
				if err != nil {
					nackErr := ec.subscription.Nack(
						"projection failed",
						esdb.NackActionRetry,
						activeSub.EventAppeared.Event,
					)
					if nackErr != nil {
						log.Fatal(nackErr)
					}
				}
				ec.subscription.Ack(activeSub.EventAppeared.Event)
			}
		}

		if activeSub.SubscriptionDropped != nil {
			break
		}
	}
}

// getAggregateIDFromStreamID is self-dscribing.
func getAggregateIDFromStreamID(streamID string) (string, error) {
	split := strings.Split(streamID, "_")
	aggregateID := split[len(split)-1]
	_, err := uuid.Parse(aggregateID)
	if err != nil {
		return "", err
	}
	return aggregateID, nil
}
