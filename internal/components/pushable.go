package components

import (
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type Pushable struct {
	Mass float32
}

func NewPushable(mass float32) *Pushable {
	return &Pushable{Mass: mass}
}

func (c *Pushable) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypePushable))
	w.WriteFloat("Mass", c.Mass, true)

	return nil
}

func (c *Pushable) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypePushable.Name)
	r.ConsumeLine()

	return nil
}
