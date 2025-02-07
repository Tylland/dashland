package components

import "github.com/tylland/dashland/game/core"

type VelocityComponent struct {
	BlockVector core.BlockVector
	Vector      core.Vector
}

func (v *VelocityComponent) IsFalling() bool {
	return v.BlockVector.Y > 0
}
