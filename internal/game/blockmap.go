package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type BlockMap struct {
	common.MapSize
	blockTextures *rl.Texture2D
	blocks        []*ecs.Entity
}

func NewBlockMap(mapSize common.MapSize, blockTextures *rl.Texture2D) *BlockMap {
	return &BlockMap{MapSize: mapSize, blockTextures: blockTextures, blocks: make([]*ecs.Entity, mapSize.Count())}
}

func (bm *BlockMap) GetBlockTextures() *rl.Texture2D {
	return bm.blockTextures
}

func (bm *BlockMap) CheckBlockAtPosition(blockType ecs.EntityType, position common.BlockPosition) bool {
	if position.X < 0 || position.X >= bm.Width || position.Y < 0 || position.Y >= bm.Height {
		return false
	}

	block := bm.blocks[position.Y*bm.Width+position.X]
	return block != nil && block.Type == blockType
}

func (bm *BlockMap) CheckNeighbourTypes(blockType ecs.EntityType, position common.BlockPosition) (neighbours [9]bool) {

	neighbours[OverLeft] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockLeft, BlockOver))
	neighbours[OverCenter] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockCenter, BlockOver))
	neighbours[OverRight] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockRight, BlockOver))

	neighbours[MiddleLeft] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockLeft, BlockMiddle))
	neighbours[MiddleCenter] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockCenter, BlockMiddle))
	neighbours[MiddleRight] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockRight, BlockMiddle))

	neighbours[UnderLeft] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockLeft, BlockUnder))
	neighbours[UnderCenter] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockCenter, BlockUnder))
	neighbours[UnderRight] = bm.CheckBlockAtPosition(blockType, position.Offset(BlockRight, BlockUnder))

	return
}

func (bm *BlockMap) PrintBlockMap() {
	fmt.Println("BlockMap")

	for y := 0; y < bm.Height; y++ {
		row := ""

		for x := 0; x < bm.Width; x++ {
			if row != "" {
				row += ", "
			}

			block := bm.blocks[y*bm.Width+x]
			if block == nil {
				row += "-"
			} else {
				row += fmt.Sprintf("%d", int(block.Type))
			}
		}

		fmt.Println(row)
	}
}

// SetBlock stores block at pos and returns the block previously stored there.
// It returns nil without changing the map if block is nil or pos is out of bounds.
func (bm *BlockMap) SetBlock(block *ecs.Entity, pos common.BlockPosition) *ecs.Entity {

	if block == nil || pos.X < 0 || pos.X >= bm.Width || pos.Y < 0 || pos.Y >= bm.Height {
		return nil
	}

	index := pos.Y*bm.Width + pos.X
	existing := bm.blocks[index]
	bm.blocks[index] = block

	return existing
}

func (bm *BlockMap) GetBlock(x int, y int) (*ecs.Entity, bool) {

	if x < 0 || x >= bm.Width || y < 0 || y >= bm.Height {
		return nil, false
	}

	block := bm.blocks[y*bm.Width+x]
	return block, block != nil
}

// GetBlockAtPosition returns the block entity at position and whether a block
// is present there. It returns (nil, false) when position is outside the map
// or the cell is empty.
func (bm *BlockMap) GetBlockAtPosition(position common.BlockPosition) (*ecs.Entity, bool) {

	if position.X < 0 || position.X >= bm.Width || position.Y < 0 || position.Y >= bm.Height {
		return nil, false
	}

	block := bm.blocks[position.Y*bm.Width+position.X]

	return block, block != nil
}

func (bm *BlockMap) GetNearbyBlocks(position *common.BlockPosition) []*ecs.Entity {
	// Get entity's block position
	blockX := position.X
	blockY := position.Y

	// Check surrounding blocks (3x3 grid around entity)
	var nearbyBlocks []*ecs.Entity

	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			checkX := blockX + x
			checkY := blockY + y

			// Check world bounds
			if checkX >= 0 && checkX < bm.Width &&
				checkY >= 0 && checkY < bm.Height {
				if block, ok := bm.GetBlock(checkX, checkY); ok && block != nil {
					nearbyBlocks = append(nearbyBlocks, block)
				}
			}
		}
	}

	return nearbyBlocks
}

// func (bm *BlockMap) Update(deltaTime float32) {
// 	for i := len(bm.blocks) - 1; i >= 0; i-- {
// 		if bm.blocks[i] != nil {
// 			UpdateBlock(bm.blocks[i], bm)
// 		}
// 	}
// }

func (bm *BlockMap) IsBlocked(pos common.BlockPosition, collider *components.ColliderComponent) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= bm.Width || pos.Y >= bm.Height {
		return true
	}

	if block, ok := bm.GetBlockAtPosition(pos); ok {
		blockCollider := ecs.GetComponent[components.ColliderComponent](block)
		if blockCollider == nil {
			return false
		}
		obstructs, _ := blockCollider.Collide(collider)
		return obstructs
	}

	return true
}

// func (bm *BlockMap) Render() {
// 	for index, block := range bm.blocks {
// 		x := index % int(bm.Width)
// 		y := index / int(bm.Width)

// 		rl.DrawTextureRec(bm.blockTextures, rl.NewRectangle(float32(block.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(x)*bm.blockWidth, float32(y)*bm.blockHeight), rl.White)
// 	}

// 	for _, block := range bm.blocks {
// 		rl.DrawTextureRec(bm.blockTextures, rl.NewRectangle(float32(block.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(block.position.X)*bm.blockWidth, float32(block.position.Y)*bm.blockHeight), rl.White)
// 	}
// }
