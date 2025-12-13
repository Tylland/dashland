package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/game/core"
)

type MapSize struct {
	width       int
	height      int
	blockWidth  float32
	blockHeight float32
}

type BlockMap struct {
	MapSize
	blockTextures *rl.Texture2D
	blocks        []*Block
}

func NewBlockMap(mapSize MapSize, blockTextures *rl.Texture2D) *BlockMap {
	return &BlockMap{MapSize: mapSize, blockTextures: blockTextures}
}

func (bm *BlockMap) InitBlocks(world *World, tiles []*tiled.LayerTile) {
	bm.blocks = make([]*Block, len(tiles))

	for index, tile := range tiles {
		bm.blocks[index] = NewBlock(world, BlockType(tile.ID), index%bm.width, index/bm.width)
	}
}

func (bm *BlockMap) CheckBlockAtPosition(blockType BlockType, position core.BlockPosition) bool {
	if position.X < 0 || position.X >= bm.width || position.Y < 0 || position.Y >= bm.height {
		return false
	}

	return bm.blocks[position.Y*bm.width+position.X].blockType == blockType
}

func (bm *BlockMap) CheckNeighbourTypes(blockType BlockType, position core.BlockPosition) (neighbours [9]bool) {

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

	for y := 0; y < bm.height; y++ {
		row := ""

		for x := 0; x < bm.width; x++ {
			if row != "" {
				row += ", "
			}

			row += fmt.Sprintf("%d", int(bm.blocks[y*bm.width+x].blockType))
		}

		fmt.Println(row)
	}
}

func (bm *BlockMap) SetBlock(block *Block, pos core.BlockPosition) {

	if pos.X < 0 || pos.X >= bm.width || pos.Y < 0 || pos.Y >= bm.height {
		return
	}

	block.position = pos

	bm.blocks[pos.Y*bm.width+pos.X] = block
}

func (bm *BlockMap) GetBlock(x int, y int) (*Block, bool) {

	if x < 0 || x >= bm.width || y < 0 || y >= bm.height {
		return &Block{blockType: Unknown}, false
	}

	return bm.blocks[y*bm.width+x], true
}

func (bm *BlockMap) GetBlockAtPosition(position core.BlockPosition) (*Block, bool) {

	if position.X < 0 || position.X >= bm.width || position.Y < 0 || position.Y >= bm.height {
		return &Block{blockType: Unknown}, false
	}

	return bm.blocks[position.Y*bm.width+position.X], true
}

func (bm *BlockMap) update(deltaTime float32) {

	for i := len(bm.blocks) - 1; i >= 0; i-- {
		bm.blocks[i].update(deltaTime)
	}
}

func (bm *BlockMap) render() {
	// for index, block := range bm.blocks {
	// 	x := index % int(bm.width)
	// 	y := index / int(bm.width)

	// 	rl.DrawTextureRec(bm.blockTextures, rl.NewRectangle(float32(block.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(x)*bm.blockWidth, float32(y)*bm.blockHeight), rl.White)
	// }

	// for _, block := range bm.blocks {
	// 	rl.DrawTextureRec(bm.blockTextures, rl.NewRectangle(float32(block.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(block.position.X)*bm.blockWidth, float32(block.position.Y)*bm.blockHeight), rl.White)
	// }
}
