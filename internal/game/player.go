package game

import (
	"math"

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

type Player struct {
	stage    *Stage
	Position *components.PositionComponent
	Velocity *components.VelocityComponent
	common.BoxBody
	IsDead   bool
	pickaxe  bool
	Movement common.Movement
}

func NewPlayer(stage *Stage) *Player {
	return &Player{stage: stage, Velocity: components.NewVelocityComponentZero()}
}

func NewPlayerEntity(world *ecs.World, stage *Stage, position common.BlockPosition) (*ecs.Entity, *ecs.Components) {
	player := ecs.NewEntity(ecs.EntityID("player"), EntityPlayer)

	comps := ecs.NewComponents()
	comps.AddComponent(components.NewCharacteristicsComponent(characteristics.IsPlayer | characteristics.CanHoldGravity))
	comps.AddComponent(components.NewInputComponent())
	comps.AddComponent(components.NewPositionComponent(position, stage.GetPosition(position)))
	comps.AddComponent(components.NewVelocityComponentZero())
	comps.AddComponent(components.NewSpriteComponent(common.NewSprite(stage.entityTextures, stage.BlockWidth, stage.BlockHeight, float32(EntityPlayer)*stage.BlockWidth, 0, 1, 0)))

	return player, comps
}

func (p *Player) InitPosition(blockPosition common.BlockPosition) {
	position := p.stage.GetPosition(blockPosition)

	p.Position = components.NewPositionComponent(blockPosition, position)

	p.BoxBody = common.BoxBody{Width: 32, Height: 32}

	//p.targetBlockPosition = blockPosition
	//p.Position = p.game.world.GetPosition(pos)
	//	p.SetPosition(position)
}

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

func (p *Player) Hurt(entity *ecs.Entity) {
	p.IsDead = true

	// position := p.Position.CurrentBlockPosition

	// for y := -1; y <= 1; y++ {
	// 	for x := -1; x <= 1; x++ {
	// 		diamondPosition := position.Offset(x, y)
	// 		if !p.blocks.CheckBlockAtPosition(game.Bedrock, diamondPosition) {
	// 			world.SetBlock(NewBlock(world, Void, diamondPosition.X, diamondPosition.Y), diamondPosition)
	// 			entity := NewDiamond(world, diamondPosition, p.blocks.GetPosition(diamondPosition))
	// 			world.SetEntity(entity, diamondPosition)
	// 		}
	// 	}
	// }

}
