package ecs

type EntityCategory uint16
type EntityType uint16
type EntityID string

type Entity struct {
	ID   EntityID
	Type EntityType
}

func NewEntity(id EntityID, entityType EntityType) *Entity {
	return &Entity{
		ID:   id,
		Type: entityType,
	}
}
