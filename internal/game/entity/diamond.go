package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

func NewDiamond(header *ecs.EntityHeader, size common.MapSize, blockPosition common.BlockPosition, position rl.Vector2) *ecs.Entity {
	entity := ecs.NewEntity(header)

	entity.AddComponent(components.NewPositionComponent(blockPosition, position))
	entity.AddComponent(components.NewCharacteristicsComponent(characteristics.CanFall | characteristics.CanHoldGravity | characteristics.GravityRollOff | characteristics.Destructable))
	entity.AddComponent(components.NewBlockStep(position))
	entity.AddComponent(components.NewSpriteComponent(common.NewSprite(assets.LoadTexture("entities"), size.BlockWidth, size.BlockHeight, float32(types.EntityDiamond)*size.BlockWidth, 0, 0)))
	entity.AddComponent(components.NewColliderComponent(collider.LayerCollectable, collider.LayerAll&(^collider.LayerPlayer|collider.LayerEnemy), collider.LayerPlayer|collider.LayerGround|collider.LayerEnemy))
	entity.AddComponent(components.NewCollectableComponent("diamond", 1))
	entity.AddComponent(components.NewDamage(1))

	return entity
}
