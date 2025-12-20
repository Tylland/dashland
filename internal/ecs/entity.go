package ecs

type EntityCategory uint16
type EntityType uint16
type EntityID string

type Entity struct {
	ID   EntityID
	Type EntityType
	//	Characteristics characteristics.Characteristics
	// Characteristic Component
	// Position       Component
	// Velocity       Component
	// Sprite         Component
	// Collision      Component
	// Collectable    Component
	// WallWalker     Component
}

func NewEntity(id EntityID, entityType EntityType) *Entity {
	return &Entity{
		ID:   id,
		Type: entityType,
	}
}
