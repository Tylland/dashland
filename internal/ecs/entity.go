package ecs

type EntityCategory uint16
type EntityType uint16
type EntityId string

type Entity struct {
	Id   EntityId
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

func NewEntity(id EntityId, entityType EntityType) *Entity {
	return &Entity{
		Id:   id,
		Type: entityType,
	}
}
