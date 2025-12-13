package game

import (
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/core"
)

type EntityCategory uint16

const (
	EntityCategoryObject EntityCategory = 0
	EntityCategoryEnemy  EntityCategory = 100
)

type EntityType uint16

const (
	EntityPlayer  EntityType = iota
	EntityDiamond            = 4
	EntityBoulder            = 5
	EntityEnemy              = 101
)

type Entity struct {
	Id              core.EntityId
	Type            EntityType
	Characteristics characteristics.Characteristics
	Position        *components.PositionComponent
	Velocity        *components.VelocityComponent
	Sprite          *components.SpriteComponent
	Collision       *components.CollisionComponent
	Collectable     *components.CollectableComponent
	WallWalker      *components.WallWalkerComponent
}

func (e *Entity) HasCharacteristic(characteristic characteristics.Characteristics) bool {
	return e.Characteristics&characteristic == characteristic
}
