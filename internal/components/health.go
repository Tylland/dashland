package components

import "github.com/tylland/dashland/internal/ecs"

type Health struct {
	Points int
}

func NewHealth(points int) *Health {
	return &Health{Points: points}
}

func (c *Health) WriteTo(w ecs.Writer) error {
	w.StartSection(ecs.ComponentName(c))
	w.WriteInt("Points", c.Points, true)

	return nil
}

func (c *Health) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel(ecs.ComponentName(c))
	r.ConsumeLine()

	return nil
}
