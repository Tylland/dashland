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

func NewMagicWall(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity | characteristics.IsMagicWall | characteristics.GravityRollOff))
	entity.AddComponent(components.NewColliderComponent(collider.LayerWall, collider.LayerAll, collider.LayerNone))

	animations := map[string]components.Animation{
		"idle":   {BaseX: 0, BaseY: 0, FrameCount: 1, FrameDuration: 1, Loop: true},
		"active": {BaseX: 0, BaseY: 0, FrameCount: 4, FrameDuration: 0.100, Loop: true},
	}

	entity.AddComponent(components.NewAnimationComponent(animations, "idle"))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("magicwall"), size.BlockWidth, size.BlockHeight, 0, 0, 0)))

	area := rl.NewRectangle(position.X, position.Y, size.BlockWidth, size.BlockHeight)
	entity.AddComponent(components.NewParticleComponent("mist", area, 1, 0.2, 32, 64, 0.20, rl.NewColor(235, 205, 205, 255)))

	return entity
}
