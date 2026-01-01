package common

import (
	"fmt"
	"strconv"
	"strings"
)

type BlockPosition struct {
	X int
	Y int
}

func NewBlockPosition(x int, y int) BlockPosition {
	return BlockPosition{X: x, Y: y}
}

func NewBlockPositionFromIndex(index, width int) BlockPosition {
	return BlockPosition{X: index % width, Y: index / width}
}

func ParseBlockPosition(str string) (BlockPosition, error) {
	pos := BlockPosition{}
	parts := strings.Split(strings.TrimSpace(str), ",")

	if len(parts) != 2 {
		return pos, fmt.Errorf("Invalid format %s", str)
	}

	x, err := strconv.Atoi(parts[0])

	if err != nil {
		return pos, err
	}

	y, err := strconv.Atoi(parts[1])

	if err != nil {
		return pos, err
	}

	return NewBlockPosition(x, y), nil
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

func (bp BlockPosition) IsZero() bool {
	return bp.X == 0 && bp.Y == 0
}
