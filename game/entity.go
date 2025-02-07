package game

import (
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/core"
)

type Entity struct {
	Id          core.EntityId
	Type        BlockType
	Behavior    BlockBehavior
	Position    *components.PositionComponent
	Velocity    *components.VelocityComponent
	Sprite      *components.SpriteComponent
	Collision   *components.CollisionComponent
	Collectable *components.CollectableComponent
}
