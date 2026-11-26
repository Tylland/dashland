package ecs

import (
	"fmt"
	"slices"
	"strings"
)

type ComponentType struct {
	Name    string
	Version int
}

type Component interface {
	//Name() string
}

type Components struct {
	index   map[string]Component
	ordered []Component
}

func NewComponents() *Components {
	return &Components{
		index: make(map[string]Component),
	}
}

func ComponentName(component Component) string {
	return strings.TrimPrefix(strings.TrimSuffix(strings.ToLower(fmt.Sprintf("%T", component)), "component"), "*components.")
}

func (c *Components) AddComponent(comp Component) {

	name := ComponentName(comp)
	c.index[name] = comp
	c.ordered = append(c.ordered, comp)
}

func (c *Components) GetComponent(name string) (Component, bool) {
	comp, ok := c.index[name]
	return comp, ok
}

func (c *Components) RemoveComponent(comp Component) {
	name := ComponentName(comp)
	delete(c.index, name)

	if i := slices.Index(c.ordered, comp); i != -1 {
		c.ordered = slices.Delete(c.ordered, i, i+1)
	}
}

func GetComponent[T Component](entity *Entity) *T {
	name := ComponentName(new(T))
	commp, ok := entity.Components.GetComponent(name)

	if ok {
		return commp.(*T)
	}

	return nil
}

func (c *Components) Ordered() []Component {
	return c.ordered
}

func RemoveComponent[T Component](entity *Entity) {
	entity.Components.RemoveComponent(new(T))
}
