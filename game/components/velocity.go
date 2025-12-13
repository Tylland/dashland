package components

import "github.com/tylland/dashland/game/core"

type VelocityComponent struct {
	BlockVector core.BlockVector
	Vector      core.Vector
}

func NewVelocityComponent(blockVector core.BlockVector, vector core.Vector) *VelocityComponent {
	return &VelocityComponent{BlockVector: blockVector, Vector: vector}
}

func NewVelocityComponentZero() *VelocityComponent {
	return NewVelocityComponent(core.NewBlockVector(0, 0), core.NewVector(0, 0))
}

func (v *VelocityComponent) IsFalling() bool {
	return v.BlockVector.Y > 0
}
