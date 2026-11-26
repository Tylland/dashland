package components

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type Characteristic struct {
	characteristics common.Characteristics
}

func NewCharacteristics(char common.Characteristics) *Characteristic {
	return &Characteristic{characteristics: char}
}

func (c *Characteristic) Has(characteristic common.Characteristics) bool {
	return c.characteristics&characteristic == characteristic
}

func (c *Characteristic) HasNot(characteristic common.Characteristics) bool {
	return c.characteristics&^characteristic == 0
}

func (c *Characteristic) Add(characteristic common.Characteristics) {
	c.characteristics = c.characteristics | characteristic
}

func (c *Characteristic) Remove(characteristic common.Characteristics) {
	c.characteristics = c.characteristics &^ characteristic
}

func (c *Characteristic) Value() common.Characteristics {
	return c.characteristics
}

func (c *Characteristic) Set(value common.Characteristics) {
	c.characteristics = value
}

func (c *Characteristic) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeCharacteristics))
	w.WriteEnum("Characteristics", c.characteristics, true, "")

	return nil
}

func (c *Characteristic) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeCharacteristics.Name)

	if enum, changed, err := r.ReadEnum(); err == nil {
		if changed {
			c.characteristics = common.Characteristics(enum)
		}
	}

	return nil
}
