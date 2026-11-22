package components

import (
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/ecs"
)

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

func (c *CharacteristicComponent) WriteTo(w ecs.Writer) error {
	w.StartSection(ecs.ComponentName(c))
	w.WriteEnum("Characteristics", c.characteristics, false, "")

	return nil
}

func (c *CharacteristicComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel(ecs.ComponentName(c))

	if enum, changed, err := r.ReadEnum(); err == nil {
		if changed {
			c.characteristics = characteristics.Characteristics(enum)
		}
	}

	return nil
}
