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

func NewExitDoor(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2, targetStage string, targetPosition common.BlockPosition) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("entities"), size.BlockWidth, size.BlockHeight, float32(header.Type)*size.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerBedrock, collider.LayerAll, collider.LayerNone))
	entity.AddComponent(components.NewDoorComponent(targetStage, targetPosition, components.DoorHidden))

	return entity
}
