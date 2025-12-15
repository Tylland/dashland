package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
)

type PositionComponent struct {
	PreviousBlockPosition common.BlockPosition
	CurrentBlockPosition  common.BlockPosition
	TargetBlockPosition   common.BlockPosition
	rl.Vector2
}

func NewPositionComponent(blockPosition common.BlockPosition, position rl.Vector2) *PositionComponent {
	return &PositionComponent{
		PreviousBlockPosition: blockPosition,
		CurrentBlockPosition:  blockPosition,
		TargetBlockPosition:   blockPosition,
		Vector2:               position,
	}
}

func (p *PositionComponent) SetTarget(offset common.BlockVector) {
	p.TargetBlockPosition = p.CurrentBlockPosition.Add(offset)
}

// check if has target position
func (p *PositionComponent) HasTarget() bool {
	return !p.TargetBlockPosition.IsSame(p.CurrentBlockPosition)
}

// Cancel target position
func (p *PositionComponent) CancelTarget() {
	p.TargetBlockPosition = p.CurrentBlockPosition
}

// Commit the target block position to the current block position
func (p *PositionComponent) UseTarget() {
	p.PreviousBlockPosition = p.CurrentBlockPosition
	p.CurrentBlockPosition = p.TargetBlockPosition
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
