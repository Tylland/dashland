package components

import (
	"github.com/tylland/dashland/game/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PositionComponent struct {
	PreviousBlockPosition core.BlockPosition
	CurrentBlockPosition  core.BlockPosition
	TargetBlockPosition   core.BlockPosition
	rl.Vector2
}

func NewPositionComponent(blockPosition core.BlockPosition, position rl.Vector2) *PositionComponent {
	return &PositionComponent{
		PreviousBlockPosition: blockPosition,
		CurrentBlockPosition:  blockPosition,
		TargetBlockPosition:   blockPosition,
		Vector2:               position,
	}
}

func (p *PositionComponent) SetTarget(offset core.BlockVector) {
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
func (p *PositionComponent) Update(position core.BlockPosition) {
	p.PreviousBlockPosition = p.CurrentBlockPosition
	p.CurrentBlockPosition = position
}

// Rollback rolls back the position to the previous block position
func (p *PositionComponent) Rollback() {
	p.CurrentBlockPosition = p.PreviousBlockPosition
}
