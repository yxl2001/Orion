package core

import "sync"

type Event struct {
	Type string
	Data any
}

type EmitFunc func(Event)

type EventBus struct {
	mu          sync.RWMutex
	subscribers []chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (b *EventBus) Subscribe() <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 128)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *EventBus) Emit(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
