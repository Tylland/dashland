package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

func NewBoulder(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPosition(blockPosition, position))
	entity.AddComponent(components.NewCharacteristics(common.CanFall | common.CanHoldGravity | common.GravityRollOff | common.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSprite(common.NewSprite(assets.LoadTexture("entities"), size.BlockWidth, size.BlockHeight, float32(types.EntityBoulder)*size.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewCollider(collider.LayerItem, collider.LayerAll, collider.LayerPlayer|collider.LayerEnemy))
	entity.AddComponent(components.NewPushable(1))
	entity.AddComponent(components.NewDamage(1))

	return entity
}
