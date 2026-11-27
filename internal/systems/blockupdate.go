package systems

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
)

type BlockUpdateSystem struct {
	stage *game.Stage
}

func NewBlockUpdateSystem(stage *game.Stage) *BlockUpdateSystem {
	return &BlockUpdateSystem{stage: stage}
}

func (s *BlockUpdateSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		if entity.Type == types.EntityVoid || entity.Type == types.EntitySoil {
			if position, block, ok := entity.ComponentTuple[components.Position, components.BlockCorners](); ok {
				s.updateBlock(entity, s.stage.BlockMap, position, block)
			}
		}
	}
}

// cornerIndex returns the corner texture index for a block corner.
// A strict soil corner requires all three neighboring cells to match,
// while void corners allow the diagonal neighbor to differ.
func cornerIndex(neighbors [9]bool, corner game.CornerPosition, strict bool) uint8 {
	if corner == game.CornerOverLeft && neighbors[game.MiddleLeft] && (neighbors[game.OverLeft] || !strict) && neighbors[game.OverCenter] {
		return 1
	}
	if corner == game.CornerOverRight && neighbors[game.OverCenter] && (neighbors[game.OverRight] || !strict) && neighbors[game.MiddleRight] {
		return 2
	}
	if corner == game.CornerUnderRight && neighbors[game.MiddleRight] && (neighbors[game.UnderRight] || !strict) && neighbors[game.UnderCenter] {
		return 3
	}
	if corner == game.CornerUnderLeft && neighbors[game.UnderCenter] && (neighbors[game.UnderLeft] || !strict) && neighbors[game.MiddleLeft] {
		return 4
	}
	return 0
}

func (s *BlockUpdateSystem) updateBlock(entity *ecs.Entity, blockMap *game.BlockMap, position *components.Position, block *components.BlockCorners) {
	switch entity.Type {
	case types.EntitySoil:
		neighbors := blockMap.CheckNeighbourTypes(types.EntityVoid, position.CurrentBlockPosition)
		block.Corners[game.CornerOverLeft] = cornerIndex(neighbors, game.CornerOverLeft, true)
		block.Corners[game.CornerOverRight] = cornerIndex(neighbors, game.CornerOverRight, true)
		block.Corners[game.CornerUnderRight] = cornerIndex(neighbors, game.CornerUnderRight, true)
		block.Corners[game.CornerUnderLeft] = cornerIndex(neighbors, game.CornerUnderLeft, true)
	case types.EntityVoid:
		neighbors := blockMap.CheckNeighbourTypes(types.EntitySoil, position.CurrentBlockPosition)
		block.Corners[game.CornerOverLeft] = cornerIndex(neighbors, game.CornerOverLeft, false)
		block.Corners[game.CornerOverRight] = cornerIndex(neighbors, game.CornerOverRight, false)
		block.Corners[game.CornerUnderRight] = cornerIndex(neighbors, game.CornerUnderRight, false)
		block.Corners[game.CornerUnderLeft] = cornerIndex(neighbors, game.CornerUnderLeft, false)
	}
}
