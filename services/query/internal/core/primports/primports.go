package primports

import "github.com/jrlangford/go-eventsourcing-cqrs/lib/query"

// QueryRunner wraps the Run method.
type QueryRunner interface {
	Run(message query.QueryMessage) (interface{}, error)
}
