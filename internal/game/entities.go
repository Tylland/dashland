package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/collider"
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
	case EntityExplosion:
		return NewExplosion(world, stage, blockPosition, position), nil
	case EntityMagicWall:
		return NewMagicWall(world, stage, blockPosition, position), nil
	// case EntityDoor:
	// 	return NewDoor(world, stage, blockPosition, position), nil
	case EntityFirefly:
		return NewFirefly(world, stage, blockPosition, position), nil
	case EntityButterfly:
		return NewButterfly(world, stage, blockPosition, position), nil
	default:
		return nil, fmt.Errorf("%v (%d) is unknown EntityType", entityType, int(entityType))

	}
}

func LoadEntity(world *ecs.World, stage *Stage, reader ecs.Reader) (*ecs.Entity, error) {
	var blockPosition common.BlockPosition

	header := &ecs.EntityHeader{}
	header.ReadFrom(reader)

	entity, err := NewGameEntity(world, stage, header.Type, blockPosition)

	if err != nil {
		return nil, err
	}

	entity.EntityHeader = header

	if reader != nil {
		entity.ReadFrom(reader)
	}

	return entity, nil
}

func NewBoulder(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("boulder"), EntityBoulder))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityBoulder)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerItem, collider.LayerAll, collider.LayerPlayer|collider.LayerEnemy))
	entity.AddComponent(components.NewPushableComponent(1))
	entity.AddComponent(components.NewDamage(1))

	world.AddEntity(entity)

	return entity
}

func NewDiamond(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("diamond"), EntityDiamond))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDiamond)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerCollectable, collider.LayerAll&(^collider.LayerPlayer|collider.LayerEnemy), collider.LayerPlayer|collider.LayerGround|collider.LayerEnemy))
	entity.AddComponent(components.NewCollectableComponent("diamond", 1))
	entity.AddComponent(components.NewDamage(1))

	world.AddEntity(entity)

	return entity
}

func NewFirefly(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("firefly"), EntityFirefly))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.IsEnemy | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewColliderComponent(collider.LayerEnemy, collider.LayerEnemy, collider.LayerPlayer))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewHealth(1))
	entity.AddComponent(components.NewWallWalkerComponent(true))

	animations := map[string]components.Animation{
		"default": {BaseX: 0, BaseY: 0 * stage.BlockHeight, FrameCount: 4, FrameDuration: 0.200, Loop: true},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "default"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("firefly"), stage.BlockWidth, stage.BlockHeight, 0, 0, 0)))

	world.AddEntity(entity)

	return entity
}

func NewButterfly(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("butterfly"), EntityButterfly))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.IsEnemy | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewColliderComponent(collider.LayerEnemy, collider.LayerEnemy, collider.LayerPlayer))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewHealth(1))
	entity.AddComponent(components.NewWallWalkerComponent(false))

	animations := map[string]components.Animation{
		"default": {BaseX: 0, BaseY: 0 * stage.BlockHeight, FrameCount: 8, FrameDuration: 0.100, Loop: true},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "default"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("butterfly"), stage.BlockWidth, stage.BlockHeight, 0, 0, 0)))

	world.AddEntity(entity)

	return entity
}

func NewDoor(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2, targetStage string, targetPosition common.BlockPosition) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("door"), EntityDoor))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDoor)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerTrigger, collider.LayerAll&(^collider.LayerPlayer|collider.LayerEnemy), collider.LayerPlayer))
	entity.AddComponent(components.NewDoorComponent(targetStage, targetPosition, components.DoorClosed))

	area := rl.NewRectangle(position.X, position.Y, stage.BlockWidth, stage.BlockHeight)
	entity.AddComponent(components.NewParticleComponent("mist", area, 40, 0.5, 32, 64, 0.20, rl.NewColor(235, 205, 205, 255)))

	world.AddEntity(entity)

	return entity
}

func NewMagicWall(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("magicwall"), EntityMagicWall))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity | characteristics.IsMagicWall | characteristics.GravityRollOff))
	entity.AddComponent(components.NewColliderComponent(collider.LayerWall, collider.LayerAll, collider.LayerNone))

	animations := map[string]components.Animation{
		"idle":   {BaseX: 0, BaseY: 0, FrameCount: 1, FrameDuration: 1, Loop: true},
		"active": {BaseX: 0, BaseY: 0, FrameCount: 4, FrameDuration: 0.100, Loop: true},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "idle"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("magicwall"), stage.BlockWidth, stage.BlockHeight, 0, 0, 0)))

	area := rl.NewRectangle(position.X, position.Y, stage.BlockWidth, stage.BlockHeight)
	entity.AddComponent(components.NewParticleComponent("mist", area, 1, 0.2, 32, 64, 0.20, rl.NewColor(235, 205, 205, 255)))

	world.AddEntity(entity)

	return entity
}

func NewExitDoor(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2, targetStage string, targetPosition common.BlockPosition) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("exitdoor"), EntityDoor))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.blockTextures, stage.BlockWidth, stage.BlockHeight, float32(Bedrock)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerBedrock, collider.LayerAll, collider.LayerNone))
	entity.AddComponent(components.NewDoorComponent(targetStage, targetPosition, components.DoorHidden))

	world.AddEntity(entity)

	return entity
}

func NewExplosion(world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("explosion"), EntityExplosion))

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.IsEnemy))
	entity.AddComponent(components.NewColliderComponent(collider.LayerEnemy, collider.LayerEnemy, collider.LayerPlayer))
	entity.AddComponent(components.NewDamage(1))

	animations := map[string]components.Animation{
		"default": {BaseX: 0, BaseY: 0 * stage.BlockHeight, FrameCount: 6, FrameDuration: 0.100, Loop: false},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "default"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("explosion"), stage.BlockWidth, stage.BlockHeight, 0, 0, 0)))
	entity.AddComponent(components.NewLifetime(0.6))
	world.AddEntity(entity)

	return entity
}

func NewFlash(world *ecs.World) *ecs.Entity {
	entity := ecs.NewEntity(ecs.NewEntityHeader(NewEntityId("flash"), EntityFlash))

	entity.AddComponent(components.NewFlashComponent())
	entity.AddComponent(components.NewLifetime(0.1))

	world.AddEntity(entity)

	return entity
}
