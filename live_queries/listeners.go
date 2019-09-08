package live_queries

import (
	"fmt"
	"log"
	"time"
)

type Listeners struct {
	listeners      map[int]*listenerWrapper
	nextListenerID int
}

type listenerWrapper struct {
	name     string
	listener *Listener
}

type Listener struct {
	Insert func(r Row) error
	Update func(before Row, after Row) error
	Delete func(r Row) error
}

type ListenerID int

func NewListeners() *Listeners {
	return &Listeners{
		listeners: map[int]*listenerWrapper{},
	}
}

func (l *Listeners) Process(evt *Event) {
	for _, listener := range l.listeners {
		start := time.Now()
		err := listener.listener.run(evt)
		dur := time.Since(start)
		if err != nil {
			log.Printf("err running listener %v: %v (%v)", listener.name, err, dur)
		} else {
			log.Printf("ran listener %v (%v)", listener.name, dur)
		}
	}
}

func (l *Listener) run(evt *Event) error {
	if evt.Payload.After != nil && evt.Payload.Before != nil {
		if l.Update != nil {
			if err := l.Update(evt.Payload.Before, evt.Payload.After); err != nil {
				return fmt.Errorf("error from update listener: %v", err)
			}
		}
		return nil
	}
	if evt.Payload.After != nil {
		if l.Insert != nil {
			if err := l.Insert(evt.Payload.After); err != nil {
				return fmt.Errorf("error from insert listener: %v", err)
			}
		}
		return nil
	}
	if evt.Payload.Before != nil {
		if l.Delete != nil {
			if err := l.Delete(evt.Payload.Before); err != nil {
				return fmt.Errorf("error from delete listener: %v", err)
			}
		}
		return nil
	}
	return nil
}

func (l *Listeners) AddListener(name string, list *Listener) ListenerID {
	lid := l.nextListenerID
	l.listeners[lid] = &listenerWrapper{
		name:     name,
		listener: list,
	}
	l.nextListenerID++
	return ListenerID(lid)
}
