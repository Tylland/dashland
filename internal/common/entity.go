package common

type ComponentType uint

type GameEntity interface {
	HasComponents(components ComponentType) bool
}

type Component interface {
}
