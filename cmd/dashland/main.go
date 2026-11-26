package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal"
	"github.com/tylland/dashland/internal/common"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Dashland")

	rl.InitAudioDevice()

	game := internal.NewGame(screenWidth, screenHeight)
	defer game.Unload()

	err := game.LoadStage("start", common.NewBlockPosition(0, 0))

	if err != nil {
		panic(err)
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		game.Update(deltaTime)
	}

	rl.CloseAudioDevice()
}
