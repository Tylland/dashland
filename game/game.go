package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const mapPath = "maps/start.tmx" // Path to your Tiled Map.

type DashlandGame struct {
	Screen
	Camera   Camera
	Player   Player
	BlockMap BlockMap
}

func NewGame(screenWidth int, screenHeight int) *DashlandGame {
	game := DashlandGame{Screen: Screen{Width: screenWidth, Height: screenHeight}}
	game.init()

	return &game
}

func (g *DashlandGame) init() {

	blockMap, err := loadBlockMapFromFile(mapPath)

	if err != nil {
		return
	}

	g.BlockMap = blockMap
	g.Player = Player{game: g, lastBlockPosition: BlockPosition{27, 2}, targetBlockPosition: BlockPosition{27, 2}}
	g.Camera = NewSmoothFollowCamera(&g.Screen, &g.Player)
}

func (g *DashlandGame) Update(deltaTime float32) {
	fmt.Println(deltaTime)

	g.BlockMap.Update(deltaTime)
	g.Player.Update(deltaTime)
	g.Camera.Update(deltaTime)
}

func (g *DashlandGame) Render() {
	// Draw
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(g.Camera.GetCamera())

	g.BlockMap.Render()
	g.Player.Render()

	rl.EndMode2D()

	rl.EndDrawing()
}
