package events

import (
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"fmt"

	githubeventbus "github.com/asaskevich/EventBus"
)

type InMemoryEventBus struct {
	bus githubeventbus.Bus
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{bus: githubeventbus.New()}
}

func (b *InMemoryEventBus) Publish(event shareddomain.DomainEvent) {
	b.bus.Publish(event.EventName(), event)
}

func (b *InMemoryEventBus) Subscribe(
	eventName string,
	handler func(shareddomain.DomainEvent),
) error {
	return b.bus.Subscribe(eventName, func(event any) {
		domainEvent, ok := event.(shareddomain.DomainEvent)
		if !ok {
			panic(fmt.Sprintf("unexpected event payload for %s", eventName))
		}

		handler(domainEvent)
	})
}
