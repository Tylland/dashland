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

type EntityCounters struct {
	Counters map[string]int
}

func NewEntityCounters() *EntityCounters {
	return &EntityCounters{Counters: make(map[string]int)}
}

func (ec *EntityCounters) Increment(prefix string) int {
	ec.Counters[prefix]++

	return ec.Counters[prefix]
}

func (ec *EntityCounters) Reset() {
	ec.Counters = make(map[string]int)
}

func (ec *EntityCounters) WriteTo(w ecs.Writer) error {
	w.StartSection("counters")

	w.WriteInt("Count", len(ec.Counters), false)

	for k, v := range ec.Counters {
		w.WriteString("Key", k, false)
		w.WriteInt("Value", v, false)
	}

	return nil
}

func (ec *EntityCounters) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel("counters")

	ec.Reset()

	count, _, err := r.ReadInt()

	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		key, _, _ := r.ReadString()
		val, _, _ := r.ReadInt()
		ec.Counters[key] = val
	}

	return nil
}

var EntityIdCounters *EntityCounters = NewEntityCounters()

func NewEntityId(prefix string) ecs.EntityID {
	count := EntityIdCounters.Increment(prefix)

	return ecs.EntityID(fmt.Sprintf("%s-%d", prefix, count))
}

func NewEntityHeader(entityType ecs.EntityType) *ecs.EntityHeader {
	prefix := EntityTypeName(entityType)
	count := EntityIdCounters.Increment(prefix)

	return ecs.NewEntityHeader(ecs.EntityID(fmt.Sprintf("%s-%d", prefix, count)), entityType, EntityCategory(entityType))
}

func NewEntity(header *ecs.EntityHeader, world *ecs.World, stage *Stage, bpos common.BlockPosition) (*ecs.Entity, error) {
	pos := stage.GetPosition(bpos)

	switch header.Type {
	case EntityVoid:
		return NewBlock(header, stage, bpos, pos, characteristics.Destructable, collider.LayerNone, collider.LayerNone, collider.LayerNone)
	case EntitySoil:
		return NewBlock(header, stage, bpos, pos, characteristics.CanHoldGravity|characteristics.Destructable, collider.LayerGround, collider.LayerEnemy, collider.LayerAll)
	case EntityWall:
		return NewBlock(header, stage, bpos, pos, characteristics.CanHoldGravity|characteristics.GravityRollOff|characteristics.Destructable, collider.LayerWall, collider.LayerAll, collider.LayerNone)
	case EntityBedrock:
		return NewBlock(header, stage, bpos, pos, characteristics.CanHoldGravity, collider.LayerBedrock, collider.LayerAll, collider.LayerNone)
	case EntityDiamond:
		return NewDiamond(header, world, stage, bpos, pos), nil
	case EntityBoulder:
		return NewBoulder(header, world, stage, bpos, pos), nil
	case EntityExplosion:
		return NewExplosion(header, world, stage, bpos, pos), nil
	case EntityMagicWall:
		return NewMagicWall(header, world, stage, bpos, pos), nil
	case EntityExitDoor:
		return NewExitDoor(header, world, stage, bpos, pos, "", bpos), nil
	case EntityFirefly:
		return NewFirefly(header, world, stage, bpos, pos), nil
	case EntityButterfly:
		return NewButterfly(header, world, stage, bpos, pos), nil
	default:
		return nil, fmt.Errorf("%q (%d) is unknown EntityType", EntityTypeName(header.Type), int(header.Type))
	}
}

func NewEntityAndHeader(entityType ecs.EntityType, world *ecs.World, stage *Stage, blockPosition common.BlockPosition) (*ecs.Entity, error) {
	header := NewEntityHeader(entityType)
	return NewEntity(header, world, stage, blockPosition)
}

func NewBoulder(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityBoulder)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerItem, collider.LayerAll, collider.LayerPlayer|collider.LayerEnemy))
	entity.AddComponent(components.NewPushableComponent(1))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewIconComponent(stage.entityTextures, rl.NewRectangle(float32(EntityBoulder)*stage.BlockWidth, 0, stage.BlockWidth, stage.BlockHeight)))

	world.AddEntity(entity)

	return entity
}

func NewDiamond(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDiamond)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerCollectable, collider.LayerAll&(^collider.LayerPlayer|collider.LayerEnemy), collider.LayerPlayer|collider.LayerGround|collider.LayerEnemy))
	entity.AddComponent(components.NewCollectableComponent("diamond", 1))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewIconComponent(stage.entityTextures, rl.NewRectangle(float32(EntityDiamond)*stage.BlockWidth, 0, stage.BlockWidth, stage.BlockHeight)))

	world.AddEntity(entity)

	return entity
}

func NewFirefly(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

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

func NewButterfly(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

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

// func NewDoor(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
// 	entity := ecs.NewEntity(header)

// 	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
// 	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
// 	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityDoor)*stage.BlockWidth, 0, 0)))
// 	entity.AddComponent(components.NewColliderComponent(collider.LayerTrigger, collider.LayerAll&(^collider.LayerPlayer|collider.LayerEnemy), collider.LayerPlayer))
// 	entity.AddComponent(&components.DoorComponent{})

// 	area := rl.NewRectangle(position.X, position.Y, stage.BlockWidth, stage.BlockHeight)
// 	entity.AddComponent(components.NewParticleComponent("mist", area, 40, 0.5, 32, 64, 0.20, rl.NewColor(235, 205, 205, 255)))

// 	world.AddEntity(entity)

// 	return entity
// }

func NewMagicWall(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

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

func NewExitDoor(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2, targetStage string, targetPosition common.BlockPosition) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("entities"), stage.BlockWidth, stage.BlockHeight, float32(header.Type)*stage.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerBedrock, collider.LayerAll, collider.LayerNone))
	entity.AddComponent(components.NewDoorComponent(targetStage, targetPosition, components.DoorHidden))

	world.AddEntity(entity)

	return entity
}

func NewExplosion(header *ecs.EntityHeader, world *ecs.World, stage *Stage, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

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
	entity := ecs.NewEntity(NewEntityHeader(EntityFlash))

	entity.AddComponent(components.NewFlashComponent())
	entity.AddComponent(components.NewLifetime(0.1))

	world.AddEntity(entity)

	return entity
}

func NewBlock(header *ecs.EntityHeader, stage *Stage, bpos common.BlockPosition, pos rl.Vector2, character characteristics.Characteristics, layer collider.CollisionLayer, obstruct collider.CollisionLayer, interact collider.CollisionLayer) (*ecs.Entity, error) {
	sprite := common.NewSprite(assets.LoadTexture("blocks"), stage.BlockWidth, stage.BlockHeight, float32(header.Type)*stage.BlockWidth, 0, 0)

	entity := ecs.NewEntity(header)
	entity.AddComponent(components.NewPositionComponent(bpos, pos))
	entity.AddComponent(components.NewSpriteComponent(sprite))
	entity.AddComponent(components.NewCharacteristicsComponent(character))
	entity.AddComponent(components.NewColliderComponent(layer, obstruct, interact))
	entity.AddComponent(&components.BlockComponent{})

	return entity, nil
}
