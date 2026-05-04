package sharedsse

import (
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"sync"
)

type Message struct {
	Data      any
	EventName string
}

type Hub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan Message]struct{}
}

func NewHub() *Hub {
	return &Hub{subscribers: map[string]map[chan Message]struct{}{}}
}

func (h *Hub) Broadcast(userID string, message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for subscriber := range h.subscribers[userID] {
		select {
		case subscriber <- message:
		default:
		}
	}
}

func (h *Hub) Subscribe(userID string) (<-chan Message, func()) {
	channel := make(chan Message, 8)

	h.mu.Lock()
	if h.subscribers[userID] == nil {
		h.subscribers[userID] = map[chan Message]struct{}{}
	}
	h.subscribers[userID][channel] = struct{}{}
	h.mu.Unlock()

	return channel, func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		delete(h.subscribers[userID], channel)
		if len(h.subscribers[userID]) == 0 {
			delete(h.subscribers, userID)
		}
		close(channel)
	}
}

type Connector struct {
	eventBus shareddomain.EventBus
	events   []string
	hub      *Hub
}

func NewConnector(eventBus shareddomain.EventBus, hub *Hub, events []string) (*Connector, error) {
	connector := &Connector{eventBus: eventBus, events: events, hub: hub}
	for _, eventName := range events {
		if err := eventBus.Subscribe(eventName, func(event shareddomain.DomainEvent) {
			hub.Broadcast(event.UserID(), Message{Data: event, EventName: event.EventName()})
		}); err != nil {
			return nil, err
		}
	}

	return connector, nil
}

func (c *Connector) ConfiguredEvents() []string {
	cloned := make([]string, len(c.events))
	copy(cloned, c.events)
	return cloned
}

func (c *Connector) Stream(writer stdhttp.ResponseWriter, request *stdhttp.Request, userID string) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	flusher, ok := writer.(stdhttp.Flusher)
	if !ok {
		writer.WriteHeader(stdhttp.StatusInternalServerError)
		return
	}

	channel, unsubscribe := c.hub.Subscribe(userID)
	defer unsubscribe()

	writer.WriteHeader(stdhttp.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-request.Context().Done():
			return
		case message, ok := <-channel:
			if !ok {
				return
			}

			payload, err := json.Marshal(message.Data)
			if err != nil {
				continue
			}

			_, _ = fmt.Fprintf(writer, "event: %s\n", message.EventName)
			_, _ = fmt.Fprintf(writer, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}
