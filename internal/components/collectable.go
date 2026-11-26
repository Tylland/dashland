package components

import (
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type CollectableType uint32

const (
	CollectableNone    CollectableType = 0
	CollectableDiamond CollectableType = 1
)

type Collectable struct {
	Name      string
	Amount    int
	Collected bool
}

func NewCollectable(name string, amount int) *Collectable {
	return &Collectable{Name: name, Amount: amount, Collected: false}
}

func (c *Collectable) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeCollectable))
	w.WriteString("Name", c.Name, false)
	w.WriteInt("Amount", c.Amount, false)

	return nil
}

func (c *Collectable) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeCollectable.Name)
	r.ConsumeLine()

	ecs.WithChangedValue(r.ReadInt, func(amount int) {
		c.Amount = amount
	})

	return nil
}
