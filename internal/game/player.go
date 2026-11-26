package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type Vector struct {
	X float32
	Y float32
}

func (v Vector) GetLength() float32 {
	return float32(math.Sqrt(float64(v.X)*float64(v.X) + float64(v.Y)*float64(v.Y)))
}

func (v Vector) Normalized() Vector {
	length := v.GetLength()

	return Vector{X: v.X / length, Y: v.Y / length}
}

func NewPlayerEntity(world *ecs.World, stage *Stage, position common.BlockPosition, texture *rl.Texture2D) *ecs.Entity {
	player := ecs.NewEntity(NewEntityHeader(types.EntityPlayer))

	player.AddComponent(components.NewCharacteristics(common.CanHoldGravity))
	player.AddComponent(components.NewInput())
	player.AddComponent(components.NewPosition(position, stage.GetPosition(position)))
	player.AddComponent(components.NewBlockStep(stage.GetPosition(position)))
	player.AddComponent(components.NewCollider(collider.LayerPlayer, collider.LayerNone, collider.LayerAll))
	player.AddComponent(components.NewSprite(common.NewSprite(texture, stage.BlockWidth, stage.BlockHeight, float32(1)*stage.BlockWidth, 0, 0)))

	// Example animations for the player: idle and walk.
	animations := map[string]common.Animation{
		"idle":        {BaseX: 0, BaseY: 0 * stage.BlockHeight, FrameCount: 1, FrameDuration: 5, Loop: true},
		"stand_left":  {BaseX: 0, BaseY: 0 * stage.BlockHeight, FrameCount: 1, FrameDuration: 5, Loop: true},
		"stand_up":    {BaseX: 1 * stage.BlockWidth, BaseY: 0 * stage.BlockHeight, FrameCount: 1, FrameDuration: 5, Loop: true},
		"stand_right": {BaseX: 2 * stage.BlockWidth, BaseY: 0 * stage.BlockHeight, FrameCount: 1, FrameDuration: 5, Loop: true},
		"stand_down":  {BaseX: 3 * stage.BlockWidth, BaseY: 0 * stage.BlockHeight, FrameCount: 1, FrameDuration: 5, Loop: true},
		"walk_left":   {BaseX: 0, BaseY: 1 * stage.BlockHeight, FrameCount: 4, FrameDuration: 0.100, Loop: true},
		"walk_up":     {BaseX: 0, BaseY: 2 * stage.BlockHeight, FrameCount: 4, FrameDuration: 0.100, Loop: true},
		"walk_right":  {BaseX: 0, BaseY: 3 * stage.BlockHeight, FrameCount: 4, FrameDuration: 0.100, Loop: true},
		"walk_down":   {BaseX: 0, BaseY: 4 * stage.BlockHeight, FrameCount: 4, FrameDuration: 0.100, Loop: true},
	}

	player.AddComponent(components.NewAnimation(animations, "idle"))
	player.AddComponent(components.NewInventory())
	player.AddComponent(components.NewHealth(1))

	return player
}
