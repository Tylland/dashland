package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func NewFirefly(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPosition(blockPosition, position))
	entity.AddComponent(components.NewCharacteristics(common.IsEnemy | common.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewCollider(collider.LayerEnemy, collider.LayerEnemy, collider.LayerPlayer))
	entity.AddComponent(components.NewDamage(1))
	entity.AddComponent(components.NewHealth(1))
	entity.AddComponent(components.NewWallWalker(true))

	animations := map[string]common.Animation{
		"default": {BaseX: 0, BaseY: 0 * size.BlockHeight, FrameCount: 4, FrameDuration: 0.200, Loop: true},
	}

	entity.AddComponent(components.NewAnimation(animations, "default"))
	entity.AddComponent(components.NewSprite(common.NewSprite(assets.LoadTexture("firefly"), size.BlockWidth, size.BlockHeight, 0, 0, 0)))

	return entity
}
