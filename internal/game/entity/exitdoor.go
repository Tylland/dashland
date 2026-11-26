package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func NewExitDoor(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2, targetStage string, targetPosition common.BlockPosition) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPosition(blockPosition, position))
	entity.AddComponent(components.NewCharacteristics(common.CanHoldGravity))
	entity.AddComponent(components.NewSprite(common.NewSprite(assets.LoadTexture("entities"), size.BlockWidth, size.BlockHeight, float32(header.Type)*size.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewCollider(collider.LayerBedrock, collider.LayerAll, collider.LayerNone))
	entity.AddComponent(components.NewDoor(targetStage, targetPosition, components.DoorHidden))

	return entity
}
