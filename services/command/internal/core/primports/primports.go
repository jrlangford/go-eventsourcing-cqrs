package primports

import es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"

// CommandExecutor wraps the Execute method.
type CommandExecutor interface {
	Execute(es.CommandMessage) error
}
