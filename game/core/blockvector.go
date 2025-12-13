package core

type BlockVector struct {
	X int
	Y int
}

func NewBlockVector(x int, y int) BlockVector {
	return BlockVector{X: x, Y: y}
}

func (v *BlockVector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v *BlockVector) Clear() {
	v.X = 0
	v.Y = 0
}

func (v *BlockVector) Offset(deltaX int, deltaY int) BlockVector {
	return BlockVector{X: v.X + deltaX, Y: v.Y + deltaY}
}

func (v *BlockVector) Add(vector BlockVector) BlockVector {
	return BlockVector{X: v.X + vector.X, Y: v.Y + vector.Y}
}

func (v *BlockVector) Subtract(vector BlockVector) BlockVector {
	return BlockVector{X: v.X - vector.X, Y: v.Y - vector.Y}
}

func (v BlockVector) TurnLeft() BlockVector {
	return BlockVector{
		X: v.Y,
		Y: -v.X,
	}
}

func (v BlockVector) TurnRight() BlockVector {
	return BlockVector{
		X: -v.Y,
		Y: v.X,
	}
}

func (v BlockVector) Reverse() BlockVector {
	return BlockVector{
		X: -v.X,
		Y: -v.Y,
	}
}

// Example usage:
// Up (0,-1) turns left to Left (-1,0)
// Left (-1,0) turns left to Down (0,1)
// Down (0,1) turns left to Right (1,0)
// Right (1,0) turns left to Up (0,-1)
