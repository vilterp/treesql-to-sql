package live_queries

import "log"

type Listeners struct {
	listeners      map[int]*Listener
	nextListenerID int
}

type Listener struct {
	Insert func(r Row) error
	Update func(before Row, after Row) error
	Delete func(r Row) error
}

type ListenerID int

func NewListeners() *Listeners {
	return &Listeners{
		listeners: map[int]*Listener{},
	}
}

func (l *Listeners) Process(evt *Event) {
	for _, listener := range l.listeners {
		if evt.Payload.After != nil && evt.Payload.Before != nil {
			if listener.Update != nil {
				if err := listener.Update(evt.Payload.Before, evt.Payload.After); err != nil {
					log.Println("error from update listener:", err)
				}
			}
			continue
		}
		if evt.Payload.After != nil {
			if listener.Insert != nil {
				if err := listener.Insert(evt.Payload.After); err != nil {
					log.Println("error from insert listener:", err)
				}
			}
			continue
		}
		if evt.Payload.Before != nil {
			if listener.Delete != nil {
				if err := listener.Delete(evt.Payload.Before); err != nil {
					log.Println("error from delete listener:", err)
				}
			}
			continue
		}
	}
}

func (l *Listeners) AddListener(list *Listener) ListenerID {
	lid := l.nextListenerID
	l.listeners[lid] = list
	l.nextListenerID++
	return ListenerID(lid)
}
