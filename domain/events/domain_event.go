package events

import "time"

// DomainEvent represents a generic domain event.
type DomainEvent interface {
	EventType() string
	OccurredOn() time.Time
}
