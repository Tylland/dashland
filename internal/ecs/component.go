package ecs

import (
	"fmt"
	"strings"
)

type Component interface {
	//Name() string
}

type Components struct {
	components map[string]Component
}

func NewComponents() *Components {
	return &Components{
		components: make(map[string]Component),
	}
}

func ComponentName(component Component) string {
	return strings.TrimPrefix(strings.TrimSuffix(strings.ToLower(fmt.Sprintf("%T", component)), "component"), "*components.")
}

func (c *Components) AddComponent(component Component) {
	name := ComponentName(component)
	c.components[name] = component
}

func (c *Components) GetComponent(name string) (Component, bool) {
	component, ok := c.components[name]
	return component, ok
}

func (c *Components) RemoveComponent(component Component) {
	name := ComponentName(component)
	delete(c.components, name)
}

func GetComponent[T Component](components *Components) *T {
	name := ComponentName(new(T))
	commp, ok := components.GetComponent(name)

	if ok {
		return commp.(*T)
	}

	return nil
}
