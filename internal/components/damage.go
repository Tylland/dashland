package components

import (
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type Damage struct {
	Points int
}

func NewDamage(points int) *Damage {
	return &Damage{Points: points}
}

func (c *Damage) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeDamage))

	w.WriteInt("Points", c.Points, true)

	return nil
}

func (c *Damage) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeDamage.Name)
	r.ConsumeLine()

	return nil
}
