package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/core"
)

var entityId core.EntityId = 0

func NewEntityId() core.EntityId {
	entityId++

	return entityId
}

func NewGameEntity(world *World, entityType EntityType, x int, y int) (*Entity, error) {
	blockPosition := core.BlockPosition{X: x, Y: y}
	position := world.GetPosition(blockPosition)

	switch entityType {
	case EntityDiamond:
		return NewDiamond(world, NewEntityId(), blockPosition, position), nil
	case EntityBoulder:
		return NewBoulder(world, NewEntityId(), blockPosition, position), nil
	case EntityEnemy:
		return NewEnemy(world, NewEntityId(), blockPosition, position), nil
	default:
		return nil, fmt.Errorf("%v (%d) is unknown EntityType", entityType, int(entityType))

	}
}

func NewBoulder(world *World, id core.EntityId, blockPosition core.BlockPosition, position rl.Vector2) *Entity {
	return &Entity{
		Id:              id,
		Type:            EntityBoulder,
		Characteristics: characteristics.CanFall | characteristics.Pushable | characteristics.EnemyObstacle,
		Position:        components.NewPositionComponent(blockPosition, position),
		Velocity:        components.NewVelocityComponentZero(),
		Sprite:          components.NewSpriteComponent(core.NewSprite(world.entityTextures, world.blockWidth, world.blockHeight, float32(EntityBoulder)*world.blockWidth, 0, 1, 0)),
		Collision:       components.NewCollisionComponent(world.blockWidth, world.blockHeight, components.LayerAll),
	}
}

func NewDiamond(world *World, id core.EntityId, blockPosition core.BlockPosition, position rl.Vector2) *Entity {
	return &Entity{
		Id:              id,
		Type:            EntityDiamond,
		Characteristics: characteristics.CanFall | characteristics.Collectable | characteristics.EnemyObstacle,
		Position:        components.NewPositionComponent(blockPosition, position),
		Velocity:        components.NewVelocityComponentZero(),
		Sprite:          components.NewSpriteComponent(core.NewSprite(world.entityTextures, world.blockWidth, world.blockHeight, float32(EntityDiamond)*world.blockWidth, 0, 1, 0)),
		Collision:       components.NewCollisionComponent(world.blockWidth, world.blockHeight, components.LayerAll),

		Collectable: components.NewCollectableComponent(components.CollectableDiamond, 1),
	}
}

func NewEnemy(world *World, id core.EntityId, blockPosition core.BlockPosition, position rl.Vector2) *Entity {
	return &Entity{
		Id:              id,
		Type:            EntityEnemy,
		Characteristics: characteristics.IsEnemy,
		Position:        components.NewPositionComponent(blockPosition, position),
		Velocity:        components.NewVelocityComponentZero(),
		Sprite:          components.NewSpriteComponent(core.NewSprite(world.enemyTextures, world.blockWidth, world.blockHeight, 0, 0, 4, 0)),
		Collision:       components.NewCollisionComponent(world.blockWidth, world.blockHeight, components.LayerAll),
		WallWalker:      &components.WallWalkerComponent{DefaultDirection: core.BlockVector{X: 1, Y: 0}},
	}
}
