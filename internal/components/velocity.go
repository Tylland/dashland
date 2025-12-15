package components

import "github.com/tylland/dashland/internal/common"

type VelocityComponent struct {
	BlockVector common.BlockVector
	Vector      common.Vector
}

func NewVelocityComponent(blockVector common.BlockVector, vector common.Vector) *VelocityComponent {
	return &VelocityComponent{BlockVector: blockVector, Vector: vector}
}

func NewVelocityComponentZero() *VelocityComponent {
	return NewVelocityComponent(common.NewBlockVector(0, 0), common.NewVector(0, 0))
}

func (v *VelocityComponent) IsFalling() bool {
	return v.BlockVector.Y > 0
}
