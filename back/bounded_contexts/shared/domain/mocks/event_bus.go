package mocks

import shareddomain "brandtoonapi/bounded_contexts/shared/domain"

type EventBusMock struct {
	PublishFunc       func(event shareddomain.DomainEvent)
	PublishedEvents   []shareddomain.DomainEvent
	SubscribeFunc     func(eventName string, handler func(shareddomain.DomainEvent)) error
	SubscribedToNames []string
}

func (m *EventBusMock) Publish(event shareddomain.DomainEvent) {
	m.PublishedEvents = append(m.PublishedEvents, event)
	if m.PublishFunc != nil {
		m.PublishFunc(event)
	}
}

func (m *EventBusMock) Subscribe(eventName string, handler func(shareddomain.DomainEvent)) error {
	m.SubscribedToNames = append(m.SubscribedToNames, eventName)
	if m.SubscribeFunc == nil {
		return nil
	}

	return m.SubscribeFunc(eventName, handler)
}
