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

func NewVoidBlock(header *ecs.EntityHeader, size common.MapSize, bpos common.BlockPosition, pos rl.Vector2) (*ecs.Entity, error) {
	sprite := common.NewSprite(assets.LoadTexture("void_soil"), size.BlockWidth, size.BlockHeight, float32(header.Type)*size.BlockWidth, 0, 0)

	entity := ecs.NewEntity(header)
	entity.AddComponent(components.NewPositionComponent(bpos, pos))
	entity.AddComponent(components.NewSpriteComponent(sprite))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.Destructable))
	entity.AddComponent(components.NewColliderComponent(collider.LayerNone, collider.LayerNone, collider.LayerNone))
	entity.AddComponent(&components.BlockComponent{HasCorners: true})

	return entity, nil
}
