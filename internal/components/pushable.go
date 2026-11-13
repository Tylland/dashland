package components

import "github.com/tylland/dashland/internal/ecs"

type PushableComponent struct {
	Mass float32
}

func NewPushableComponent(mass float32) *PushableComponent {
	return &PushableComponent{Mass: mass}
}

func (c *PushableComponent) WriteTo(w ecs.Writer) error {
	w.WriteLabel(ecs.ComponentName(c))
	w.WriteFloat(c.Mass, true)

	return nil
}

func (c *PushableComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeLine()
	r.ConsumeLine()

	return nil
}
