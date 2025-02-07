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
