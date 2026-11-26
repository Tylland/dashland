package components

import (
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

const CharacteristicComponentVersion = 1

type CharacteristicComponent struct {
	characteristics characteristics.Characteristics
}

func NewCharacteristicsComponent(char characteristics.Characteristics) *CharacteristicComponent {
	return &CharacteristicComponent{characteristics: char}
}

func (c *CharacteristicComponent) Has(characteristic characteristics.Characteristics) bool {
	return c.characteristics&characteristic == characteristic
}

func (c *CharacteristicComponent) HasNot(characteristic characteristics.Characteristics) bool {
	return c.characteristics&^characteristic == 0
}

func (c *CharacteristicComponent) Add(characteristic characteristics.Characteristics) {
	c.characteristics = c.characteristics | characteristic
}

func (c *CharacteristicComponent) Remove(characteristic characteristics.Characteristics) {
	c.characteristics = c.characteristics &^ characteristic
}

func (c *CharacteristicComponent) Value() characteristics.Characteristics {
	return c.characteristics
}

func (c *CharacteristicComponent) Set(value characteristics.Characteristics) {
	c.characteristics = value
}

func (c *CharacteristicComponent) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeCharacteristics))
	w.WriteEnum("Characteristics", c.characteristics, false, "")

	return nil
}

func (c *CharacteristicComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeCharacteristics.Name)

	if enum, changed, err := r.ReadEnum(); err == nil {
		if changed {
			c.characteristics = characteristics.Characteristics(enum)
		}
	}

	return nil
}
