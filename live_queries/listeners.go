package live_queries

type Listeners struct {
	listeners      map[int]*Listener
	nextListenerID int
}

type Listener struct {
	Insert func(r Row)
	Update func(before Row, after Row)
	Delete func(r Row)
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
				listener.Update(evt.Payload.Before, evt.Payload.After)
			}
			break
		}
		if evt.Payload.After != nil {
			if listener.Insert != nil {
				listener.Insert(evt.Payload.After)
			}
			break
		}
		if evt.Payload.Before != nil {
			if listener.Delete != nil {
				listener.Delete(evt.Payload.Before)
			}
			break
		}
	}
}

func (l *Listeners) AddListener(list *Listener) ListenerID {
	lid := l.nextListenerID
	l.listeners[lid] = list
	l.nextListenerID++
	return ListenerID(lid)
}
