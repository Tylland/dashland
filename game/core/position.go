package core

type BlockPosition struct {
	X int
	Y int
}

// type Position struct {
// 	X float32
// 	Y float32
// }

func (bp BlockPosition) Offset(deltaX int, deltaY int) BlockPosition {
	return BlockPosition{X: bp.X + deltaX, Y: bp.Y + deltaY}
}

func (bp BlockPosition) Add(vector BlockVector) BlockPosition {
	return BlockPosition{X: bp.X + int(vector.X), Y: bp.Y + int(vector.Y)}
}

func (bp BlockPosition) Subtract(position BlockPosition) BlockPosition {
	return BlockPosition{X: bp.X - position.X, Y: bp.Y - position.Y}
}

func (bp BlockPosition) IsSame(other BlockPosition) bool {
	return bp.X == other.X && bp.Y == other.Y
}
