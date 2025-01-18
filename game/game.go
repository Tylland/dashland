package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const mapPath = "maps/start.tmx" // Path to your Tiled Map.

type DashlandGame struct {
	Screen
	Camera Camera
	player *Player
	//	BlockMap BlockMap
	world *world
}

func NewGame(screenWidth int, screenHeight int) *DashlandGame {
	game := DashlandGame{Screen: Screen{Width: screenWidth, Height: screenHeight}}
	game.init()

	return &game
}

func (g *DashlandGame) init() {

	g.world = NewWorld()
	g.world.initFromFile(mapPath)

	g.player = &Player{game: g, blockPosition: BlockPosition{27, 2}, targetBlockPosition: BlockPosition{27, 2}}
	g.world.initPlayer(g.player)

	g.Camera = NewSmoothFollowCamera(&g.Screen, g.player)
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
