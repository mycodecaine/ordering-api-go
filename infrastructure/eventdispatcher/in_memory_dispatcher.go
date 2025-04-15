package eventdispatcher

// infrastructure/eventdispatcher/simple_dispatcher.go
import (
	"ORDERING-API/domain/events"
	"log"
	"sync"
)

type InMemoryDispatcher struct {
	handlers map[string][]events.EventHandler
	mu       sync.RWMutex
}

func NewSimpleDispatcher() *InMemoryDispatcher {
	return &InMemoryDispatcher{
		handlers: make(map[string][]events.EventHandler),
	}
}

func (d *InMemoryDispatcher) Register(eventType string, handler events.EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *InMemoryDispatcher) Dispatch(domainEvents []events.DomainEvent) {
	for _, event := range domainEvents {
		eventType := event.EventType()
		log.Printf("Dispatch event: %s", eventType)
		d.mu.RLock()
		handlers := d.handlers[eventType]
		d.mu.RUnlock()

		if len(handlers) == 0 {
			log.Printf("No handlers registered for event type: %s", eventType)
			continue
		}

		for _, handler := range handlers {
			go handler.Handle(event) // async handling
		}
	}
}
