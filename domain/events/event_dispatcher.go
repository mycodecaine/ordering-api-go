// domain/events/event_dispatcher.go
package events

type EventHandler interface {
	Handle(event DomainEvent)
}

type EventDispatcher interface {
	Register(eventType string, handler EventHandler)
	Dispatch(events []DomainEvent)
}
