package components

import "github.com/tylland/dashland/internal/characteristics"

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
