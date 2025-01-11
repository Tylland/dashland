package game

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

type BlockMap struct {
	width         int
	height        int
	blockWidth    float32
	blockHeight   float32
	blockTextures rl.Texture2D
	//	blockTypes    []uint32
	blocks []Block
}

// nocommit
// func (bm *BlockMap) getBlockTypes(tiles []*tiled.LayerTile) []uint32 {
// 	blockTypes := make([]uint32, len(tiles))

// 	for index, tile := range tiles {
// 		blockTypes[index] = tile.ID
// 	}

// 	return blockTypes
// }

func (bm *BlockMap) createBlock(tile *tiled.LayerTile, x int, y int) Block {

	switch BlockType(tile.ID) {
	case Diamond:
		return Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}, flags: CanFall | Collectable}
	case Boulder:
		return Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}, flags: CanFall | Obstacle}
	default:
		return Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}}
	}

}

func (bm *BlockMap) createBlocks(tiles []*tiled.LayerTile) []Block {
	blocks := make([]Block, len(tiles))

	for index, tile := range tiles {
		blocks[index] = bm.createBlock(tile, index%bm.width, index/bm.width)
	}

	return blocks

}

func loadBlockMapFromFile(filepath string) (BlockMap, error) {
	var blockMap BlockMap

	tiledMap, err := tiled.LoadFile(filepath)

	if err != nil {
		return blockMap, err
	}

	relativeImagePath := tiledMap.Tilesets[0].Image.Source

	fileName := path.Join(path.Dir(mapPath), relativeImagePath)

	blockTexture := rl.LoadTexture(fileName)

	fmt.Print(blockTexture)

	blockMap = BlockMap{
		width:         tiledMap.Width,
		height:        tiledMap.Height,
		blockWidth:    float32(tiledMap.TileWidth),
		blockHeight:   float32(tiledMap.TileHeight),
		blockTextures: blockTexture,
	}

	//blockMap.blockTypes = blockMap.getBlockTypes(tiledMap.Layers[0].Tiles)
	blockMap.blocks = blockMap.createBlocks(tiledMap.Layers[0].Tiles)

	return blockMap, err
}

func (bm *BlockMap) SetBlock(block *Block, pos BlockPosition) {

	if pos.X < 0 || pos.X >= bm.width || pos.Y < 0 || pos.Y >= bm.height {
		return
	}

	block.position = pos

	bm.blocks[pos.Y*bm.width+pos.X] = *block
}

func (bm *BlockMap) GetBlock(x int, y int) (*Block, bool) {

	if x < 0 || x >= bm.width || y < 0 || y >= bm.height {
		return &Block{blockType: Unknown}, false
	}

	return &bm.blocks[y*bm.width+x], true
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
		temp := source.blockType
		bm.Print()
		source.blockType = target.blockType
		target.blockType = temp
		bm.Print()
	}

}

func (bm *BlockMap) ObstacleForPlayer(player *Player, position BlockPosition) bool {
	block, success := bm.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	return block.ObstacleForPlayer(player)
}

func (bm *BlockMap) VisitBlock(position BlockPosition) {

	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, bm.blocks[position.Y*bm.width+position.X].blockType)
	bm.blocks[position.Y*bm.width+position.X].blockType = Void
	fmt.Printf(" to %d \n", bm.blocks[position.Y*bm.width+position.X].blockType)
}

func (bm *BlockMap) Update(deltaTime float32) {

	for i := len(bm.blocks) - 1; i >= 0; i-- {
		bm.blocks[i].Update(deltaTime)
	}
}

func (bm *BlockMap) Render() {
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
