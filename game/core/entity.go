package core

type EntityId uint32
type ComponentType uint

type GameEntity interface {
	HasComponents(components ComponentType) bool
}

type Component interface {
}
