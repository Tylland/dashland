package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
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
func (c *PositionComponent) GetBlockTarget(velocity *VelocityComponent) common.BlockPosition {
	return c.CurrentBlockPosition.Add(velocity.BlockVector)
}

// Update BlockPosition and saves
func (c *PositionComponent) Update(position common.BlockPosition) {
	c.PreviousBlockPosition = c.CurrentBlockPosition
	c.CurrentBlockPosition = position
}

// Rollback rolls back the position to the previous block position
func (c *PositionComponent) Rollback() {
	c.CurrentBlockPosition = c.PreviousBlockPosition
}

func (c *PositionComponent) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypePosition))

	w.WriteInt("X", c.CurrentBlockPosition.X, false)
	w.WriteInt("Y", c.CurrentBlockPosition.Y, false)
	w.WriteFloat("x", c.Vector2.X, false)
	w.WriteFloat("y", c.Vector2.Y, false)

	return nil
}

func (c *PositionComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypePosition.Name)

	ecs.WithChangedValue(r.ReadInt, func(x int) {
		c.CurrentBlockPosition.X = x
	})

	ecs.WithChangedValue(r.ReadInt, func(y int) {
		c.CurrentBlockPosition.Y = y
	})

	ecs.WithChangedValue(r.ReadFloat, func(x float32) {
		c.Vector2.X = x
	})

	ecs.WithChangedValue(r.ReadFloat, func(y float32) {
		c.Vector2.Y = y
	})

	return nil
}
