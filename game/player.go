package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
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
	game                *DashlandGame
	blockPosition       BlockPosition
	targetBlockPosition BlockPosition
	BoxBody
	pickaxe  bool
	movement Movement
}

func NewPlayer(g *DashlandGame) *Player {
	return &Player{game: g}
}

func (p *Player) InitPosition(blockPosition BlockPosition) {
	position := p.game.world.GetPosition(blockPosition)

	p.BoxBody = BoxBody{position: position, Width: 32, Height: 32}
	p.blockPosition = blockPosition
	p.targetBlockPosition = blockPosition
	//p.Position = p.game.world.GetPosition(pos)
	//	p.SetPosition(position)
}

func (p *Player) UpdateTargetPosition() {

	if !p.movement.moving {
		if rl.IsKeyDown(rl.KeyRight) {
			p.targetBlockPosition = p.blockPosition.Offset(1, 0)
		} else if rl.IsKeyDown(rl.KeyLeft) {
			p.targetBlockPosition = p.blockPosition.Offset(-1, 0)
		} else if rl.IsKeyDown(rl.KeyDown) {
			p.targetBlockPosition = p.blockPosition.Offset(0, 1)
		} else if rl.IsKeyDown(rl.KeyUp) {
			p.targetBlockPosition = p.blockPosition.Offset(0, -1)
		}

		if p.game.world.IsObstacleForPlayer(p, p.targetBlockPosition) {
			p.targetBlockPosition = p.blockPosition
		}
	}
}

func (p *Player) update(deltaTime float32) {
	const speed float32 = 128 //pixel per sec

	p.UpdateTargetPosition()

	if !p.blockPosition.IsSame(p.targetBlockPosition) {
		p.movement = Movement{}

		targetPosition := p.game.world.GetPosition(p.targetBlockPosition)
		p.movement.Start(p.game.world.GetPosition(p.blockPosition), targetPosition, speed, nil)
		p.game.world.VisitBlock(p.targetBlockPosition)
		p.game.world.VisitObject(p, p.targetBlockPosition)

		p.blockPosition = p.targetBlockPosition

		p.SetPosition(targetPosition)
	}

	if p.movement.moving {
		p.movement.Update(deltaTime)
		p.position = p.movement.position
	}
}

func (p *Player) Body() Body {
	return &p.BoxBody
}

func (p *Player) render() {
	rl.DrawCircle(int32(p.position.X+16), int32(p.position.Y+16), 16, rl.Gold)

	rl.DrawRectangleRec(p.BoxBody.Rectangle(), rl.Red)
}

func (p *Player) Collect(co CollectableObject) {
	co.Collected()
}
