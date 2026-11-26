package components

import "github.com/tylland/dashland/internal/ecs"

type Damage struct {
	Points int
}

func NewDamage(points int) *Damage {
	return &Damage{Points: points}
}

func (c *Damage) WriteTo(w ecs.Writer) error {
	w.StartSection(ecs.ComponentName(c))
	w.WriteInt("Points", c.Points, true)

	return nil
}

func (c *Damage) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(ecs.ComponentName(c))
	r.ConsumeLine()

	return nil
}
