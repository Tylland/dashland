package components

import "github.com/tylland/dashland/internal/ecs"

type Damage struct {
	Points int
}

func NewDamage(points int) *Damage {
	return &Damage{Points: points}
}

func (c *Damage) WriteTo(w ecs.Writer) error {
	w.WriteLabel(ecs.ComponentName(c))
	w.WriteInt(c.Points, true)

	return nil
}
