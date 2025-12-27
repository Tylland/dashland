package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
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

// type Player struct {
// 	stage    *Stage
// 	Position *components.PositionComponent
// 	Velocity *components.VelocityComponent
// 	common.BoxBody
// 	IsDead   bool
// 	pickaxe  bool
// 	Movement common.Movement
// }

// func NewPlayer(stage *Stage) *Player {
// 	return &Player{stage: stage, Velocity: components.NewVelocityComponentZero()}
// }

func NewPlayerEntity(world *ecs.World, stage *Stage, position common.BlockPosition, texture *rl.Texture2D) *ecs.Entity {
	player := ecs.NewEntity(ecs.EntityID("player"), EntityPlayer)

	player.AddComponent(components.NewCharacteristicsComponent(characteristics.CanHoldGravity))
	player.AddComponent(components.NewInputComponent())
	player.AddComponent(components.NewPositionComponent(position, stage.GetPosition(position)))
	player.AddComponent(components.NewBlockStep(stage.GetPosition(position)))
	player.AddComponent(components.NewColliderComponent(LayerPlayer, LayerNone, LayerAll))
	player.AddComponent(components.NewSpriteComponent(common.NewSprite(texture, stage.BlockWidth, stage.BlockHeight, float32(1)*stage.BlockWidth, 0, 0)))

	// Example animations for the player: idle and walk.
	animations := map[string]components.Animation{
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

	player.AddComponent(components.NewAnimationComponent(animations, "idle"))
	player.AddComponent(components.NewInventory())
	player.AddComponent(components.NewHealth(1))

	return player
}

// func (p *Player) InitPosition(blockPosition common.BlockPosition) {
// 	position := p.stage.GetPosition(blockPosition)

// 	p.Position = components.NewPositionComponent(blockPosition, position)

// 	p.BoxBody = common.BoxBody{Width: 32, Height: 32}

// 	//p.targetBlockPosition = blockPosition
// 	//p.Position = p.game.world.GetPosition(pos)
// 	//	p.SetPosition(position)
// }

// func (p *Player) UpdateTargetPosition(world *ecs.World) {

// 	if !p.Movement.Moving {

// 		if rl.IsKeyDown(rl.KeyRight) {
// 			p.Position.SetTarget(common.NewBlockVector(1, 0))
// 		} else if rl.IsKeyDown(rl.KeyLeft) {
// 			p.Position.SetTarget(common.NewBlockVector(-1, 0))
// 		} else if rl.IsKeyDown(rl.KeyDown) {
// 			p.Position.SetTarget(common.NewBlockVector(0, 1))
// 		} else if rl.IsKeyDown(rl.KeyUp) {
// 			p.Position.SetTarget(common.NewBlockVector(0, -1))
// 		}

// 		if p.stage.IsObstacleForPlayer(world, p.Position.TargetBlockPosition) {
// 			p.Position.CancelTarget()
// 		}

// 	}
// }

// func (p *Player) Update(world *ecs.World, deltaTime float32) {
// 	if p.IsDead {
// 		return
// 	}

// 	p.UpdateTargetPosition(world)
// }

// func (p *Player) Render() {
// 	if p.IsDead {
// 		return
// 	}

// 	rl.DrawCircle(int32(p.Position.Vector2.X+16), int32(p.Position.Vector2.Y+16), 16, rl.Gold)

// 	//rl.DrawRectangleRec(rl.Rectangle{X: p.Position.Vector2.X, Y: p.Position.Vector2.Y, Width: p.BoxBody.Width, Height: p.BoxBody.Height}, rl.Red)
// }

// func (p *Player) Hurt(entity *ecs.Entity) {
// 	p.IsDead = true

// 	// position := p.Position.CurrentBlockPosition

// 	// for y := -1; y <= 1; y++ {
// 	// 	for x := -1; x <= 1; x++ {
// 	// 		diamondPosition := position.Offset(x, y)
// 	// 		if !p.blocks.CheckBlockAtPosition(game.Bedrock, diamondPosition) {
// 	// 			world.SetBlock(NewBlock(world, Void, diamondPosition.X, diamondPosition.Y), diamondPosition)
// 	// 			entity := NewDiamond(world, diamondPosition, p.blocks.GetPosition(diamondPosition))
// 	// 			world.SetEntity(entity, diamondPosition)
// 	// 		}
// 	// 	}
// 	// }

// }
