package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

var entityIdCounters map[string]int = make(map[string]int)

func NewEntityId(prefix string) ecs.EntityID {
	entityIdCounters[prefix]++

	return ecs.EntityID(fmt.Sprintf("%s-%d", prefix, entityIdCounters[prefix]))
}

func NewGameEntity(world *ecs.World, stage *Stage, entityType ecs.EntityType, blockPosition common.BlockPosition) (*ecs.Entity, error) {
	position := stage.GetPosition(blockPosition)

	switch entityType {
	case EntityDiamond:
		return NewDiamond(world, stage, blockPosition, position), nil
	case EntityBoulder:
		return NewBoulder(world, stage, blockPosition, position), nil
	case EntityDoor:
		return NewDoor(world, stage, blockPosition, position), nil
	case EntityEnemy:
		return NewEnemy(world, stage, blockPosition, position), nil
	default:
		return nil, fmt.Errorf("%v (%d) is unknown EntityType", entityType, int(entityType))

	}
}

func NewBoulder(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("boulder"), EntityBoulder)

	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Destructable))
	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityBoulder)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(LayerItem, LayerAll, LayerPlayer|LayerEnemy))
	entity.AddComponent(components.NewPushableComponent(1))
	entity.AddComponent(components.NewDamage(1))

	world.AddEntity(entity)

	return entity
}

func NewDiamond(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("diamond"), EntityDiamond)

	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Obstacle | characteristics.Destructable))
	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDiamond)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(LayerCollectable, LayerAll&(^LayerPlayer|LayerEnemy), LayerPlayer|LayerGround|LayerEnemy))
	entity.AddComponent(components.NewCollectableComponent(components.CollectableDiamond, 1))
	entity.AddComponent(components.NewDamage(1))

	world.AddEntity(entity)

	return entity
}

func NewEnemy(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("enemy"), EntityEnemy)

	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.IsEnemy | characteristics.Destructable))
	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewColliderComponent(LayerEnemy, LayerEnemy, LayerPlayer))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewHealth(1))
	entity.AddComponent(components.NewWallWalkerComponent(common.NewBlockVector(1, 0)))

	animations := map[string]components.Animation{
		"default": {BaseX: 0, BaseY: 0 * stage.BlockHeight, FrameCount: 4, FrameDuration: 0.200, Loop: true},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "default"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.enemyTextures, stage.BlockWidth, stage.BlockHeight, 0, 0, 0)))

	world.AddEntity(entity)

	return entity
}

func NewDoor(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("door"), EntityDoor)

	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDoor)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(LayerTrigger, LayerAll&(^LayerPlayer|LayerEnemy), LayerPlayer))

	world.AddEntity(entity)

	return entity
}

func NewDoorWithDestination(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2, targetStage string, targetPosition common.BlockPosition) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("door"), EntityDoor)

	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDoor)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(LayerTrigger, LayerNone, LayerPlayer))
	entity.AddComponent(components.NewDoorComponent(targetStage, targetPosition))

	world.AddEntity(entity)

	return entity
}
