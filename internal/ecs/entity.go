package ecs

type EntityCategory uint16
type EntityType uint16
type EntityID string

type Entity struct {
	ID         EntityID
	Type       EntityType
	Components *Components
}

func NewEntity(id EntityID, entityType EntityType) *Entity {
	return &Entity{
		ID:         id,
		Type:       entityType,
		Components: NewComponents(),
	}
}

func (e *Entity) AddComponent(c Component) {
	e.Components.AddComponent(c)
}

func (e *Entity) RemoveComponent(c Component) {
	e.Components.RemoveComponent(c)
}
