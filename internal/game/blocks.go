package game

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

const (
	BlockWidth  float32 = 32
	BlockHeight float32 = 32
)

const (
	BlockLeft   = -1
	BlockCenter = 0
	BlockRight  = 1
	BlockOver   = -1
	BlockMiddle = 0
	BlockUnder  = 1
)

const (
	OverLeft = iota
	OverCenter
	OverRight
	MiddleLeft
	MiddleCenter
	MiddleRight
	UnderLeft
	UnderCenter
	UnderRight
)

type CornerPosition uint16

const (
	CornerOverLeft CornerPosition = iota
	CornerOverRight
	CornerUnderLeft
	CornerUnderRight
)

// cornerIndex returns the corner texture index for a block corner.
// A strict soil corner requires all three neighboring cells to match,
// while void corners allow the diagonal neighbor to differ.
func cornerIndex(neighbors [9]bool, corner CornerPosition, strict bool) uint8 {
	if corner == CornerOverLeft && neighbors[MiddleLeft] && (neighbors[OverLeft] || !strict) && neighbors[OverCenter] {
		return 1
	}
	if corner == CornerOverRight && neighbors[OverCenter] && (neighbors[OverRight] || !strict) && neighbors[MiddleRight] {
		return 2
	}
	if corner == CornerUnderRight && neighbors[MiddleRight] && (neighbors[UnderRight] || !strict) && neighbors[UnderCenter] {
		return 3
	}
	if corner == CornerUnderLeft && neighbors[UnderCenter] && (neighbors[UnderLeft] || !strict) && neighbors[MiddleLeft] {
		return 4
	}
	return 0
}

func UpdateBlock(entity *ecs.Entity, blockMap *BlockMap) {
	position := ecs.GetComponent[components.PositionComponent](entity)
	block := ecs.GetComponent[components.BlockComponent](entity)
	if position == nil || block == nil {
		return
	}

	switch entity.Type {
	case types.EntitySoil:
		neighbors := blockMap.CheckNeighbourTypes(types.EntityVoid, position.CurrentBlockPosition)
		block.Corners[CornerOverLeft] = cornerIndex(neighbors, CornerOverLeft, true)
		block.Corners[CornerOverRight] = cornerIndex(neighbors, CornerOverRight, true)
		block.Corners[CornerUnderRight] = cornerIndex(neighbors, CornerUnderRight, true)
		block.Corners[CornerUnderLeft] = cornerIndex(neighbors, CornerUnderLeft, true)
	case types.EntityVoid:
		neighbors := blockMap.CheckNeighbourTypes(types.EntitySoil, position.CurrentBlockPosition)
		block.Corners[CornerOverLeft] = cornerIndex(neighbors, CornerOverLeft, false)
		block.Corners[CornerOverRight] = cornerIndex(neighbors, CornerOverRight, false)
		block.Corners[CornerUnderRight] = cornerIndex(neighbors, CornerUnderRight, false)
		block.Corners[CornerUnderLeft] = cornerIndex(neighbors, CornerUnderLeft, false)
	}
}
