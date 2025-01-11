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

type BlockPosition struct {
	X int
	Y int
}

func (bp BlockPosition) Offset(deltaX int, deltaY int) BlockPosition {
	return BlockPosition{X: bp.X + deltaX, Y: bp.Y + deltaY}
}

func (bp BlockPosition) Add(position BlockPosition) BlockPosition {
	return BlockPosition{X: bp.X + position.X, Y: bp.Y + position.Y}
}

func (bp BlockPosition) IsSame(other BlockPosition) bool {
	return bp.X == other.X && bp.Y == other.Y
}

type Player struct {
	game                *DashlandGame
	lastBlockPosition   BlockPosition
	targetBlockPosition BlockPosition
	Position            rl.Vector2
	TargetVector        rl.Vector2
	Progress            float32
	Speed               float32
	Pickaxe             bool
	IsMoving            bool
	movement            Movement
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) UpdateTargetPosition() {

	if !p.movement.moving {
		if rl.IsKeyDown(rl.KeyRight) {
			p.targetBlockPosition = p.lastBlockPosition.Offset(1, 0)
			p.TargetVector = rl.NewVector2(32, 0)
		} else if rl.IsKeyDown(rl.KeyLeft) {
			p.targetBlockPosition = p.lastBlockPosition.Offset(-1, 0)
			p.TargetVector = rl.NewVector2(-32, 0)
		} else if rl.IsKeyDown(rl.KeyDown) {
			p.targetBlockPosition = p.lastBlockPosition.Offset(0, 1)
			p.TargetVector = rl.NewVector2(0, 32)
		} else if rl.IsKeyDown(rl.KeyUp) {
			p.targetBlockPosition = p.lastBlockPosition.Offset(0, -1)
			p.TargetVector = rl.NewVector2(0, -32)
		}

		if p.game.BlockMap.ObstacleForPlayer(p, p.targetBlockPosition) {
			p.targetBlockPosition = p.lastBlockPosition
			p.TargetVector = rl.NewVector2(0, 0)
		}
	}
}

func (p *Player) Update(deltaTime float32) {
	const speed float32 = 128 //pixel per sec
	//	const blockSize float32 = 32 //pixel per sec

	p.UpdateTargetPosition()

	if !p.lastBlockPosition.IsSame(p.targetBlockPosition) {
		p.movement = Movement{}
		p.movement.Start(p.game.BlockMap.getPosition(p.lastBlockPosition), p.game.BlockMap.getPosition(p.targetBlockPosition), speed)
		p.game.BlockMap.VisitBlock(p.targetBlockPosition)
		p.lastBlockPosition = p.targetBlockPosition
	}

	if p.movement.moving {
		p.movement.Update(deltaTime)
		p.Position = p.movement.position
	}

	// if p.IsMoving {
	// 	p.Progress += speed * deltaTime

	// 	movment := rl.Vector2Scale(p.TargetVector, rl.Clamp(p.Progress, 0, blockSize)/blockSize)

	// 	p.Position = rl.NewVector2(float32(p.lastBlockPosition.X)*32+16+movment.X, float32(p.lastBlockPosition.Y)*32+16+movment.Y)

	// 	if p.Progress >= blockSize {
	// 		p.TargetVector = rl.NewVector2(0, 0)
	// 		p.Progress = 0
	// 		p.IsMoving = false
	// 	}

	// }
}

func (p *Player) Render() {
	rl.DrawCircle(int32(p.Position.X+16), int32(p.Position.Y+16), 16, rl.Red)
}
