package game

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

type BlockType uint32

const (
	Unknown BlockType = iota
	Bedrock
	Void
	Dirt
	Diamond
	Boulder
)

type BlockMap struct {
	width         int
	height        int
	blockWidth    float32
	blockHeight   float32
	blockTextures rl.Texture2D
	blockTypes    []uint32
	blocks        []Block
}

func (bm BlockMap) getBlockTypes(tiles []*tiled.LayerTile) []uint32 {
	blockTypes := make([]uint32, len(tiles))

	for index, tile := range tiles {
		blockTypes[index] = tile.ID
	}

	return blockTypes
}

func (bm *BlockMap) createBlock(tile *tiled.LayerTile, x int, y int) Block {

	switch tile.ID {
	default:
		return Block{blockMap: bm, blockType: BlockType(tile.ID), position: BlockPosition{X: x, Y: y}}
	}

}

func (bm BlockMap) createBlocks(tiles []*tiled.LayerTile) []Block {
	blocks := make([]Block, len(tiles))

	for index, tile := range tiles {
		blocks[index] = bm.createBlock(tile, index%bm.width, index/bm.width)
	}

	return blocks

}

func loadBlockMapFromFile(filepath string) (BlockMap, error) {
	var blockMap BlockMap

	tiledMap, err := tiled.LoadFile(filepath)

	if err == nil {
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

	blockMap.blockTypes = blockMap.getBlockTypes(tiledMap.Layers[0].Tiles)
	blockMap.blocks = blockMap.createBlocks(tiledMap.Layers[0].Tiles)

	return blockMap, err
}

func (bm BlockMap) SetBlock(block Block, x int, y int) {

	if x < 0 || x >= bm.width || y < 0 || y >= bm.height {
		return
	}

	block.position.X = x
	block.position.Y = y

	bm.blocks[y*bm.width+x] = block
}

func (bm BlockMap) GetBlock(x int, y int) (Block, bool) {

	if x < 0 || x >= bm.width || y < 0 || y >= bm.height {
		return Block{blockType: Unknown}, false
	}

	return bm.blocks[y*bm.width+x], true
}

// func (bm BlockMap) PositionIsBlocked(position BlockPosition) bool {
// 	return bm.IsBlocking(position.X, position.Y)
// }

func (bm BlockMap) ObstacleForPlayer(player *Player, position BlockPosition) bool {
	block, success := bm.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	return block.ObstacleForPlayer(player)
}

func (bm BlockMap) VisitBlock(position BlockPosition) {
	bm.blockTypes[position.Y*bm.width+position.X] = uint32(Void)
}

func (bm BlockMap) Update(deltaTime float32) {

	for i := len(bm.blocks) - 1; i >= 0; i-- {
		bm.blocks[i].Update(deltaTime)
	}
}

func (bm BlockMap) Render() {

	for y := 0; y < bm.height; y++ {
		for x := 0; x < bm.width; x++ {
			rl.DrawTextureRec(bm.blockTextures, rl.NewRectangle(float32(bm.blockTypes[y*bm.width+x])*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(x)*bm.blockWidth, float32(y)*bm.blockHeight), rl.White)
		}
	}
}

type Block struct {
	blockMap  *BlockMap
	blockType BlockType
	position  BlockPosition
}

func NewBlock(blockMap *BlockMap, blockType BlockType, position BlockPosition) Block {
	return Block{blockMap: blockMap, blockType: blockType, position: position}
}

func (b Block) Update(deltaTime float32) {
	if b.blockType == Boulder {
		under := b.position.Move(0, -1)
		block, succes := b.blockMap.GetBlock(under.X, under.Y)

		if succes {
			if block.blockType == Void {
				b.blockMap.SetBlock(block, b.position.X, b.position.Y)
				b.blockMap.SetBlock(b, block.position.X, block.position.Y)
			}
		}
	}
}

func (b Block) ObstacleForPlayer(player *Player) bool {
	return b.blockType == Bedrock || b.blockType == Boulder
}

type IBlock interface {
	ObstacleForPlayer(player *Player) bool
}

type UnknownBlock struct {
}

func (ub UnknownBlock) ObstacleForPlayer(player *Player) bool {
	return false
}

type BoulderBlock struct {
}

func (b BoulderBlock) ObstacleForPlayer(player *Player) bool {
	return !player.Pickaxe
}
