package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func NewExplosion(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPosition(blockPosition, position))
	entity.AddComponent(components.NewCharacteristics(common.IsEnemy))
	entity.AddComponent(components.NewCollider(collider.LayerEnemy, collider.LayerEnemy, collider.LayerPlayer))
	entity.AddComponent(components.NewDamage(1))

	animations := map[string]common.Animation{
		"default": {BaseX: 0, BaseY: 0 * size.BlockHeight, FrameCount: 6, FrameDuration: 0.100, Loop: false},
	}

	entity.AddComponent(components.NewAnimation(animations, "default"))
	entity.AddComponent(components.NewSprite(common.NewSprite(assets.LoadTexture("explosion"), size.BlockWidth, size.BlockHeight, 0, 0, 0)))
	entity.AddComponent(components.NewLifetime(0.6))

	return entity
}
