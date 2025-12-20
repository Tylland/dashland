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
	case EntityEnemy:
		return NewEnemy(world, stage, blockPosition, position), nil
	default:
		return nil, fmt.Errorf("%v (%d) is unknown EntityType", entityType, int(entityType))

	}
}

func NewBoulder(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("boulder"), EntityBoulder)

	comps := ecs.NewComponents()
	comps.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.Pushable | characteristics.EnemyObstacle))
	comps.AddComponent(components.NewPositionComponent(blockPosition, position))
	comps.AddComponent(components.NewVelocityComponentZero())
	comps.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityBoulder)*stage.BlockWidth, 0, 1, 0)))
	comps.AddComponent(components.NewCollisionComponent(stage.BlockWidth, stage.BlockHeight, components.LayerAll))
	comps.AddComponent(components.NewPushableComponent(1))

	world.AddEntity(entity, comps)

	return entity
}

func NewDiamond(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("diamond"), EntityDiamond)

	comps := ecs.NewComponents()
	comps.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.Collectable | characteristics.EnemyObstacle))
	comps.AddComponent(components.NewPositionComponent(blockPosition, position))
	comps.AddComponent(components.NewVelocityComponentZero())
	comps.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDiamond)*stage.BlockWidth, 0, 1, 0)))
	comps.AddComponent(components.NewCollisionComponent(stage.BlockWidth, stage.BlockHeight, components.LayerAll))
	comps.AddComponent(components.NewCollectableComponent(components.CollectableDiamond, 1))

	world.AddEntity(entity, comps)

	return entity
}

func NewEnemy(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityId("enemy"), EntityEnemy)

	comps := ecs.NewComponents()
	comps.AddComponent(components.NewCharacteristicsComponent(characteristics.IsEnemy))
	comps.AddComponent(components.NewPositionComponent(blockPosition, position))
	comps.AddComponent(components.NewVelocityComponentZero())
	comps.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.enemyTextures, stage.BlockWidth, stage.BlockHeight, 0, 0, 4, 0)))
	comps.AddComponent(components.NewCollisionComponent(stage.BlockWidth, stage.BlockHeight, components.LayerAll))
	comps.AddComponent(components.NewWallWalkerComponent(common.NewBlockVector(1, 0)))

	world.AddEntity(entity, comps)

	return entity
}
