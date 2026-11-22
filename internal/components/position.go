package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
)

type PositionComponent struct {
	PreviousBlockPosition common.BlockPosition
	CurrentBlockPosition  common.BlockPosition
	//	TargetBlockPosition   common.BlockPosition
	rl.Vector2
}

func NewPositionComponent(blockPosition common.BlockPosition, position rl.Vector2) *PositionComponent {
	return &PositionComponent{
		PreviousBlockPosition: blockPosition,
		CurrentBlockPosition:  blockPosition,
		//		TargetBlockPosition:   blockPosition,
		Vector2: position,
	}
}

func NewPositionComponentZero() *PositionComponent {
	return &PositionComponent{
		PreviousBlockPosition: common.NewBlockPosition(0, 0),
		CurrentBlockPosition:  common.NewBlockPosition(0, 0),
		//		TargetBlockPosition:   common.NewBlockPosition(0, 0),
		Vector2: rl.NewVector2(0, 0),
	}
}

// func (p *PositionComponent) SetTarget(offset common.BlockVector) {
// 	p.TargetBlockPosition = p.CurrentBlockPosition.Add(offset)
// 	fmt.Printf("SetTarget to %v\n", p.TargetBlockPosition)
// }

// // check if has target position
// func (p *PositionComponent) HasTarget() bool {
// 	return !p.TargetBlockPosition.IsSame(p.CurrentBlockPosition)
// }

// // Cancel target position
// func (p *PositionComponent) CancelTarget() {
// 	p.TargetBlockPosition = p.CurrentBlockPosition
// 	fmt.Printf("CancelTarget to %v\n", p.TargetBlockPosition)
// }

// // Commit the target block position to the current block position
// func (p *PositionComponent) UseTarget() {
// 	p.PreviousBlockPosition = p.CurrentBlockPosition
// 	p.CurrentBlockPosition = p.TargetBlockPosition
// }

// Update BlockPosition and saves
func (p *PositionComponent) GetBlockTarget(velocity *VelocityComponent) common.BlockPosition {
	return p.CurrentBlockPosition.Add(velocity.BlockVector)
}

// Update BlockPosition and saves
func (p *PositionComponent) Update(position common.BlockPosition) {
	p.PreviousBlockPosition = p.CurrentBlockPosition
	p.CurrentBlockPosition = position
}

// Rollback rolls back the position to the previous block position
func (p *PositionComponent) Rollback() {
	p.CurrentBlockPosition = p.PreviousBlockPosition
}

func (p *PositionComponent) WriteTo(ew ecs.Writer) error {
	ew.StartSection(ecs.ComponentName(p))
	ew.WriteInt("X", p.CurrentBlockPosition.X, false)
	ew.WriteInt("Y", p.CurrentBlockPosition.Y, false)
	ew.WriteFloat("x", p.Vector2.X, false)
	ew.WriteFloat("y", p.Vector2.Y, false)

	return nil
}

func (p *PositionComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel(ecs.ComponentName(p))

	ecs.WithChangedValue(r.ReadInt, func(x int) {
		p.CurrentBlockPosition.X = x
	})

	ecs.WithChangedValue(r.ReadInt, func(y int) {
		p.CurrentBlockPosition.Y = y
	})

	ecs.WithChangedValue(r.ReadFloat, func(x float32) {
		p.Vector2.X = x
	})

	ecs.WithChangedValue(r.ReadFloat, func(y float32) {
		p.Vector2.Y = y
	})

	return nil
}
