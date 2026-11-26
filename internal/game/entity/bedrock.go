package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func NewBedrockBlock(header *ecs.EntityHeader, size common.MapSize, bpos common.BlockPosition, pos rl.Vector2) (*ecs.Entity, error) {
	sprite := common.NewSprite(assets.LoadTexture("blocks"), size.BlockWidth, size.BlockHeight, float32(header.Type)*size.BlockWidth, 0, 0)

	entity := ecs.NewEntity(header)
	entity.AddComponent(components.NewPosition(bpos, pos))
	entity.AddComponent(components.NewSprite(sprite))
	entity.AddComponent(components.NewCharacteristics(common.CanHoldGravity))
	entity.AddComponent(components.NewCollider(collider.LayerBedrock, collider.LayerAll, collider.LayerNone))
	entity.AddComponent(&components.BlockCorners{HasCorners: false})

	return entity, nil
}
