package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/core"
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
	game     *DashlandGame
	Position *components.PositionComponent
	Velocity *components.VelocityComponent
	core.BoxBody
	IsDead   bool
	pickaxe  bool
	movement Movement
}

func NewPlayer(g *DashlandGame) *Player {
	return &Player{game: g, Velocity: components.NewVelocityComponentZero()}
}

func (p *Player) InitPosition(blockPosition core.BlockPosition) {
	position := p.game.world.GetPosition(blockPosition)

	p.Position = components.NewPositionComponent(blockPosition, position)

	p.BoxBody = core.BoxBody{Width: 32, Height: 32}

	//p.targetBlockPosition = blockPosition
	//p.Position = p.game.world.GetPosition(pos)
	//	p.SetPosition(position)
}

func (p *Player) UpdateTargetPosition() {

	if !p.movement.moving {

		if rl.IsKeyDown(rl.KeyRight) {
			p.Position.SetTarget(core.NewBlockVector(1, 0))
		} else if rl.IsKeyDown(rl.KeyLeft) {
			p.Position.SetTarget(core.NewBlockVector(-1, 0))
		} else if rl.IsKeyDown(rl.KeyDown) {
			p.Position.SetTarget(core.NewBlockVector(0, 1))
		} else if rl.IsKeyDown(rl.KeyUp) {
			p.Position.SetTarget(core.NewBlockVector(0, -1))
		}

		if p.game.world.IsObstacleForPlayer(p, p.Position.TargetBlockPosition) {
			p.Position.CancelTarget()
		}

	}
}

func (p *Player) update(deltaTime float32) {
	if p.IsDead {
		return
	}

	p.UpdateTargetPosition()
}

func (p *Player) render() {
	if p.IsDead {
		return
	}

	rl.DrawCircle(int32(p.Position.Vector2.X+16), int32(p.Position.Vector2.Y+16), 16, rl.Gold)

	//rl.DrawRectangleRec(rl.Rectangle{X: p.Position.Vector2.X, Y: p.Position.Vector2.Y, Width: p.BoxBody.Width, Height: p.BoxBody.Height}, rl.Red)
}

func (p *Player) Hurt(entity *Entity) {
	p.IsDead = true

	world := p.game.world
	position := p.Position.CurrentBlockPosition

	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			diamondPosition := position.Offset(x, y)
			if !world.CheckBlockAtPosition(Bedrock, diamondPosition) {
				world.SetBlock(NewBlock(world, Void, diamondPosition.X, diamondPosition.Y), diamondPosition)
				entity := NewDiamond(world, NewEntityId(), diamondPosition, world.GetPosition(diamondPosition))
				world.SetEntity(entity, diamondPosition)
			}
		}
	}

}
