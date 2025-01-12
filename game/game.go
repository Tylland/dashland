package game

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

const mapPath = "maps/start.tmx" // Path to your Tiled Map.

type DashlandGame struct {
	Screen
	Camera Camera
	Player Player
	//	BlockMap BlockMap
	world *world
}

func NewGame(screenWidth int, screenHeight int) *DashlandGame {
	game := DashlandGame{Screen: Screen{Width: screenWidth, Height: screenHeight}}
	game.init()

	return &game
}

func (g *DashlandGame) loadBlockMapFromFile(filepath string) (BlockMap, error) {
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

	fmt.Printf("Reading tiles from layer %s \n", tiledMap.Layers[0].Name)
	blockMap.blocks = blockMap.createBlocks(g.world, tiledMap.Layers[0].Tiles)

	return blockMap, err
}

func (g *DashlandGame) init() {

	g.Player = Player{game: g, blockPosition: BlockPosition{27, 2}, targetBlockPosition: BlockPosition{27, 2}}

	g.world = NewWorld(&g.Player)

	blockMap, err := g.loadBlockMapFromFile(mapPath)

	if err != nil {
		return
	}

	g.world.blockMap = &blockMap

	g.Camera = NewSmoothFollowCamera(&g.Screen, &g.Player)
}

func (g *DashlandGame) Update(deltaTime float32) {
	//fmt.Println(deltaTime)
	g.world.update(deltaTime)

	// g.BlockMap.Update(deltaTime)
	// g.Player.update(deltaTime)

	g.Camera.Update(deltaTime)
}

func (g *DashlandGame) Render() {
	// Draw
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(g.Camera.GetCamera())

	g.world.render()

	rl.EndMode2D()

	rl.EndDrawing()
}
