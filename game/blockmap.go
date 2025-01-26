package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

type MapSize struct {
	width       int
	height      int
	blockWidth  float32
	blockHeight float32
}

type BlockMap struct {
	MapSize
	blockTextures rl.Texture2D
	blocks        []*Block
}

func (bm *BlockMap) InitBlocks(world *world, tiles []*tiled.LayerTile) {
	bm.blocks = make([]*Block, len(tiles))

	for index, tile := range tiles {
		bm.blocks[index] = NewBlock(world, BlockType(tile.ID), index%bm.width, index/bm.width)
	}
}

func (bm *BlockMap) CheckBlockAtPosition(blockType BlockType, position BlockPosition) bool {
	return bm.blocks[position.Y*bm.width+position.X].blockType == blockType
}

func (bm *BlockMap) CheckNeighbourTypes(blockType BlockType, position BlockPosition) (neighbours [9]bool) {
	x := position.X
	y := position.Y

	neighbours[OverLeft] = bm.blocks[(y+BlockOver)*bm.width+x+BlockLeft].blockType == blockType
	neighbours[OverCenter] = bm.blocks[(y+BlockOver)*bm.width+x+BlockCenter].blockType == blockType
	neighbours[OverRight] = bm.blocks[(y+BlockOver)*bm.width+x+BlockRight].blockType == blockType

	neighbours[MiddleLeft] = bm.blocks[(y+BlockMiddle)*bm.width+x+BlockLeft].blockType == blockType
	neighbours[MiddleCenter] = bm.blocks[(y+BlockMiddle)*bm.width+x+BlockCenter].blockType == blockType
	neighbours[MiddleRight] = bm.blocks[(y+BlockMiddle)*bm.width+x+BlockRight].blockType == blockType

	neighbours[UnderLeft] = bm.blocks[(y+BlockUnder)*bm.width+x+BlockLeft].blockType == blockType
	neighbours[UnderCenter] = bm.blocks[(y+BlockUnder)*bm.width+x+BlockCenter].blockType == blockType
	neighbours[UnderRight] = bm.blocks[(y+BlockUnder)*bm.width+x+BlockRight].blockType == blockType

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

func (bm *BlockMap) SwapBlock(source *Block, pos BlockPosition) {
	target, succes := bm.GetBlock(pos.X, pos.Y)

	if succes {
		tempBlockType := source.blockType
		tempBehavior := source.behavior
		bm.PrintBlockMap()
		source.blockType = target.blockType
		source.behavior = target.behavior
		target.blockType = tempBlockType
		target.behavior = tempBehavior
		bm.PrintBlockMap()
	}
}

func (bm *BlockMap) SetBlock(block *Block, pos BlockPosition) {

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
