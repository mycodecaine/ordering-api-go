package aggregate

import "ORDERING-API/domain/events"

type AggregateRoot struct {
	Events []events.DomainEvent
}

func (a *AggregateRoot) RecordEvent(event events.DomainEvent) {
	a.Events = append(a.Events, event)
}

func (a *AggregateRoot) GetEvents() []events.DomainEvent {
	return a.Events
}

func (a *AggregateRoot) ClearEvents() {
	a.Events = nil
}
