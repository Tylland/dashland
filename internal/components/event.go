package components

type Event struct {
	Name string
	Data any
}

// func NewEvent(name string, data any) *Event {
// 	return &Event{Name: name, Data: data}
// }

type EventComponent struct {
	events []*Event
}

func NewEventComponent() *EventComponent {
	return &EventComponent{events: make([]*Event, 0)}
}

func (c *EventComponent) Add(name string, data any) {
	c.events = append(c.events, &Event{Name: name, Data: data})
}

func (c *EventComponent) Events() []*Event {
	return c.events
}

func (c *EventComponent) Clear() {
	c.events = make([]*Event, 0)
}
