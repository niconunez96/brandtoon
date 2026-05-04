package shareddomain

type DomainEvent interface {
	EventName() string
	UserID() string
}

type EventBus interface {
	Publish(event DomainEvent)
	Subscribe(eventName string, handler func(DomainEvent)) error
}
