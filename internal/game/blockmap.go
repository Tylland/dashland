package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
)

type BlockMap struct {
	MapSize
	blockTextures *rl.Texture2D
	blocks        []*Block
}

func NewBlockMap(mapSize MapSize, blockTextures *rl.Texture2D) *BlockMap {
	return &BlockMap{MapSize: mapSize, blockTextures: blockTextures}
}

func (bm *BlockMap) CheckBlockAtPosition(blockType BlockType, position common.BlockPosition) bool {
	if position.X < 0 || position.X >= bm.Width || position.Y < 0 || position.Y >= bm.Height {
		return false
	}

	return bm.blocks[position.Y*bm.Width+position.X].BlockType == blockType
}

func (bm *BlockMap) CheckNeighbourTypes(blockType BlockType, position common.BlockPosition) (neighbours [9]bool) {

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

			row += fmt.Sprintf("%d", int(bm.blocks[y*bm.Width+x].BlockType))
		}

		fmt.Println(row)
	}
}

func (bm *BlockMap) SetBlock(block *Block, pos common.BlockPosition) {

	if pos.X < 0 || pos.X >= bm.Width || pos.Y < 0 || pos.Y >= bm.Height {
		return
	}

	block.Position = pos

	bm.blocks[pos.Y*bm.Width+pos.X] = block
}

func (bm *BlockMap) GetBlock(x int, y int) (*Block, bool) {

	if x < 0 || x >= bm.Width || y < 0 || y >= bm.Height {
		return &Block{BlockType: Unknown}, false
	}

	return bm.blocks[y*bm.Width+x], true
}

func (bm *BlockMap) GetBlockAtPosition(position common.BlockPosition) (*Block, bool) {

	if position.X < 0 || position.X >= bm.Width || position.Y < 0 || position.Y >= bm.Height {
		return &Block{BlockType: Unknown}, false
	}

	return bm.blocks[position.Y*bm.Width+position.X], true
}

func (bm *BlockMap) GetNearbyBlocks(position *common.BlockPosition) []*Block {
	// Get entity's block position
	blockX := position.X
	blockY := position.Y

	// Check surrounding blocks (3x3 grid around entity)
	var nearbyBlocks []*Block

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

func (bm *BlockMap) Update(deltaTime float32) {

	for i := len(bm.blocks) - 1; i >= 0; i-- {
		bm.blocks[i].Update(deltaTime)
	}
}

func (bm *BlockMap) IsBlocked(pos common.BlockPosition, collider *components.ColliderComponent) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= bm.Width || pos.Y >= bm.Height {
		return true
	}

	if block, ok := bm.GetBlockAtPosition(pos); ok {
		blocked, _ := block.Collider.Result(collider)
		return blocked
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
