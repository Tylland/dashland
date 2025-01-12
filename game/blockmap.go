package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

type BlockMap struct {
	width         int
	height        int
	blockWidth    float32
	blockHeight   float32
	blockTextures rl.Texture2D
	blocks        []*Block
}

// func (bm *BlockMap) createBlock(tile *tiled.LayerTile, x int, y int) *Block {
// 	return NewBlock(bm, BlockType(tile.ID), x, y)

// 	// switch BlockType(tile.ID) {
// 	// case Diamond:
// 	// 	return &Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}, behavior: CanFall | Collectable}
// 	// case Boulder:
// 	// 	return &Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}, behavior: CanFall | Obstacle}
// 	// default:
// 	// 	return &Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}, behavior: NoBehavior}
// 	// }

// }

func (bm *BlockMap) createBlocks(world *world, tiles []*tiled.LayerTile) []*Block {
	blocks := make([]*Block, len(tiles))

	for index, tile := range tiles {
		blocks[index] = NewBlock(world, BlockType(tile.ID), index%bm.width, index/bm.width)
	}

	return blocks

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

func (bm *BlockMap) getPosition(position BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*bm.blockWidth, float32(position.Y)*bm.blockHeight)
}

func (bm *BlockMap) CheckTypeAtPosition(blockType BlockType, position BlockPosition) bool {
	return bm.blocks[position.Y*bm.width+position.X].blockType == blockType
}

func (bm *BlockMap) SwapBlock(source *Block, pos BlockPosition) {
	target, succes := bm.GetBlock(pos.X, pos.Y)

	if succes {
		tempBlockType := source.blockType
		tempBehavior := source.behavior
		bm.Print()
		source.blockType = target.blockType
		source.behavior = target.behavior
		target.blockType = tempBlockType
		target.behavior = tempBehavior
		bm.Print()
	}

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

	for _, block := range bm.blocks {
		rl.DrawTextureRec(bm.blockTextures, rl.NewRectangle(float32(block.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(block.position.X)*bm.blockWidth, float32(block.position.Y)*bm.blockHeight), rl.White)
	}
}

func (bm *BlockMap) Print() {
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
