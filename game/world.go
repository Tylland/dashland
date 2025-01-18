package game

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

type actor interface {
	update(deltaTime float32)
	render()
}

type world struct {
	BlockMap
	player *Player
	actors []actor
}

func NewWorld() *world {
	return &world{actors: []actor{}}
}

func (w *world) initPlayer(player *Player) {
	w.player = player
	w.addActor(player)
}

func (w *world) initFromFile(filepath string) error {

	tiledMap, err := tiled.LoadFile(filepath)

	if err != nil {
		return err
	}

	relativeImagePath := tiledMap.Tilesets[0].Image.Source

	fileName := path.Join(path.Dir(mapPath), relativeImagePath)

	blockTexture := rl.LoadTexture(fileName)

	fmt.Print(blockTexture)

	w.width = tiledMap.Width

	w.height = tiledMap.Height
	w.blockWidth = float32(tiledMap.TileWidth)
	w.blockHeight = float32(tiledMap.TileHeight)
	w.blockTextures = blockTexture

	fmt.Printf("Reading tiles from layer %s \n", tiledMap.Layers[0].Name)
	w.blocks = w.createBlocks(tiledMap.Layers[0].Tiles)

	return nil
}

func (w *world) createBlocks(tiles []*tiled.LayerTile) []*Block {
	blocks := make([]*Block, len(tiles))

	for index, tile := range tiles {
		blocks[index] = NewBlock(w, BlockType(tile.ID), index%w.width, index/w.width)
	}

	return blocks

}

func (w *world) update(deltaTime float32) {
	for _, block := range w.blocks {
		block.update(deltaTime)
	}

	for _, act := range w.actors {
		act.update(deltaTime)
	}
}

func (w *world) render() {
	for _, block := range w.blocks {
		rl.DrawTextureRec(w.blockTextures, rl.NewRectangle(float32(block.blockType)*w.blockWidth, 0, w.blockWidth, w.blockHeight), rl.NewVector2(float32(block.position.X)*w.blockWidth, float32(block.position.Y)*w.blockHeight), rl.White)
	}

	for _, act := range w.actors {
		act.render()
	}
}

func (w *world) addActor(actor actor) {
	w.actors = append(w.actors, actor)
}

func (w *world) obstacleForPlayer(player *Player, position BlockPosition) bool {
	block, success := w.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	return block.ObstacleForPlayer(player)
}

func (w *world) VisitBlock(position BlockPosition) {

	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, w.blocks[position.Y*w.width+position.X].blockType)
	w.blocks[position.Y*w.width+position.X] = NewBlock(w, Void, position.X, position.Y)
	fmt.Printf(" to %d \n", w.blocks[position.Y*w.width+position.X].blockType)
}

func (w *world) checkPositionOccupied(position BlockPosition) bool {
	return !w.CheckTypeAtPosition(Void, position) || w.player.blockPosition.IsSame(position)
}

func (w *world) checkPlayerAtPosition(position BlockPosition) bool {
	return w.player.blockPosition.IsSame(position)
}
