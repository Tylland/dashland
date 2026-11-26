package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func NewButterfly(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.IsEnemy | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewColliderComponent(collider.LayerEnemy, collider.LayerEnemy, collider.LayerPlayer))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewHealth(1))
	entity.AddComponent(components.NewWallWalkerComponent(false))

	animations := map[string]components.Animation{
		"default": {BaseX: 0, BaseY: 0 * size.BlockHeight, FrameCount: 8, FrameDuration: 0.100, Loop: true},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "default"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("butterfly"), size.BlockWidth, size.BlockHeight, 0, 0, 0)))

	return entity
}
