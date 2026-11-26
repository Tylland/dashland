package components

import (
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type PushableComponent struct {
	Mass float32
}

func NewPushableComponent(mass float32) *PushableComponent {
	return &PushableComponent{Mass: mass}
}

func (c *PushableComponent) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypePushable))
	w.WriteFloat("Mass", c.Mass, true)

	return nil
}

func (c *PushableComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypePushable.Name)
	r.ConsumeLine()

	return nil
}
