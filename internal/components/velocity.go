package components

import (
	"github.com/tylland/dashland/internal/common"
)

type VelocityComponent struct {
	BlockTarget common.BlockPosition
	BlockVector common.BlockVector
	Vector      common.Vector
}

func NewVelocityComponent(blockVector common.BlockVector, vector common.Vector) *VelocityComponent {
	return &VelocityComponent{BlockVector: blockVector, Vector: vector}
}

func NewVelocityComponentZero() *VelocityComponent {
	return NewVelocityComponent(common.NewBlockVector(0, 0), common.NewVector(0, 0))
}

func (v *VelocityComponent) Clear() {
	v.BlockVector.Clear()
	v.Vector.Clear()
}

func (v *VelocityComponent) Set(bv common.BlockVector) {
	v.BlockVector = bv
}

func (v *VelocityComponent) IsMoving() bool {
	return v.BlockVector.X != 0 || v.BlockVector.Y != 0
}

func (v *VelocityComponent) IsFalling() bool {
	return v.BlockVector.Y > 0
}
