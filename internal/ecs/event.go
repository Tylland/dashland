package ecs

type Event struct {
	Name string
	Data any
}

func NewEvent(name string, data any) *Event {
	return &Event{Name: name, Data: data}
}
