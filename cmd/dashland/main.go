package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal"
	"github.com/tylland/dashland/internal/common"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Game
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

	// End Game

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		// Lookup("allocs") creates a profile similar to go test -memprofile.
		// Alternatively, use Lookup("heap") for a profile
		// that has inuse_space as the default index.
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}
