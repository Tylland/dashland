package ecs

import (
	"fmt"
	"strings"
)

type ComponentType struct {
	Name    string
	Version int
}

type Component interface {
	//Name() string
}

// type ComponentKey reflect.Type
type ComponentKey string

// func key[T Component]() ComponentKey { return reflect.TypeFor[T]() } // Go 1.22+
// func keyOf(x Component) ComponentKey { return reflect.TypeOf(x) }

func key[T Component]() ComponentKey { return ComponentKey(ComponentName(new(T))) } // Go 1.22+
func keyOf(c Component) ComponentKey { return ComponentKey(ComponentName(c)) }

func ComponentName(component Component) string {
	return strings.TrimPrefix(strings.TrimSuffix(strings.ToLower(fmt.Sprintf("%T", component)), "component"), "*components.")
}

// type Components struct {
// 	//index   map[string]Component
// 	index   map[componentKey]Component
// 	ordered []Component
// }

// func NewComponents() *Components {
// 	return &Components{
// 		index: make(map[componentKey]Component),
// 	}
// }

// func (c *Components) AddComponent(comp Component) {
// 	//name := ComponentName(comp)

// 	c.index[keyOf(comp)] = comp
// 	c.ordered = append(c.ordered, comp)
// }

// func (c *Components) GetComponent(key componentKey) (Component, bool) {
// 	comp, ok := c.index[key]
// 	return comp, ok
// }

// func (c *Components) Component[T any]() (*T, bool) {
// 	if comp, ok := c.index[key[T]()]; ok {
// 		return comp.(*T), true
// 	}

// 	return nil, false
// }

// func (c *Components) RemoveComponent(comp Component) {
// 	//name := ComponentName(comp)
// 	delete(c.index, keyOf(comp))

// 	if i := slices.Index(c.ordered, comp); i != -1 {
// 		c.ordered = slices.Delete(c.ordered, i, i+1)
// 	}
// }

// func (c *Components) Ordered() []Component {
// 	return c.ordered
// }

// func GetComponent[T Component](entity *Entity) *T {
// 	commp, ok := entity.Components.GetComponent()

// 	if ok {
// 		return commp.(*T)
// 	}

// 	return nil
// }
