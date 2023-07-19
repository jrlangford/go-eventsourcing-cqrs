package primports

// EventProjector wraps the ProjectEvent method.
type EventProjector interface {
	ProjectEvent(aggregateID string, domainEvent interface{}, version *int, metadata []byte) error
}
