package components

import (
	"github.com/tylland/dashland/internal/common"
)

type Velocity struct {
	BlockTarget common.BlockPosition
	BlockVector common.BlockVector
	Vector      common.Vector
}

func NewVelocity(blockVector common.BlockVector, vector common.Vector) *Velocity {
	return &Velocity{BlockVector: blockVector, Vector: vector}
}

func NewVelocityComponentZero() *Velocity {
	return NewVelocity(common.NewBlockVector(0, 0), common.NewVector(0, 0))
}

func (v *Velocity) Clear() {
	v.BlockVector.Clear()
	v.Vector.Clear()
}

func (v *Velocity) Set(bv common.BlockVector) {
	v.BlockVector = bv
}

func (v *Velocity) IsMoving() bool {
	return v.BlockVector.X != 0 || v.BlockVector.Y != 0
}

func (v *Velocity) IsFalling() bool {
	return v.BlockVector.Y > 0
}
