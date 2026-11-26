package components

import "github.com/tylland/dashland/internal/ecs"

type PushableComponent struct {
	Mass float32
}

func NewPushableComponent(mass float32) *PushableComponent {
	return &PushableComponent{Mass: mass}
}

func (c *PushableComponent) WriteTo(w ecs.Writer) error {
	const version = 1

	w.StartSection(ecs.ComponentName(c), version)

	w.WriteFloat("Mass", c.Mass, true)

	return nil
}

func (c *PushableComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(ecs.ComponentName(c))
	r.ConsumeLine()

	return nil
}
