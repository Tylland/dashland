package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
)

const (
	MapsRootPath   = `C:\src\dashland\assets\maps`
	StagesRootPath = `C:\src\dashland\assets\stages`
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func main() {
	mapFiles := must(filepath.Glob(MapsRootPath + "\\" + "*.tmx"))

	fmt.Println(mapFiles)

	for _, mapFile := range mapFiles {
		convertMap(mapFile)
	}
}

func convertMap(fileName string) {
	name := strings.Replace(filepath.Base(fileName), ".tmx", "", 1)

	fmt.Printf("converting map: %q\n", name)

	tiledMap := must(tiled.LoadFile(fileName))

	mapSize := common.NewMapSize(tiledMap.Width, tiledMap.Height, tiledMap.TileWidth, tiledMap.TileHeight)
	stage := game.NewStage(name, mapSize, nil, nil, nil)
	world := ecs.NewWorld()

	if tiledMap.Properties != nil {
		stage.DiamondsRequired = tiledMap.Properties.GetInt("DiamondsRequired")
		stage.DiamondPoints = tiledMap.Properties.GetInt("DiamondPoints")
		stage.DiamondBonusPoints = tiledMap.Properties.GetInt("DiamondBonusPoints")
		stage.SetTimeLimit(tiledMap.Properties.GetInt("TimeLimit"))

		game.EntityIdCounters.Reset()

		fmt.Printf("Reading blocks from layer %s \n", tiledMap.Layers[0].Name)
		game.InitBlocks(world, stage, tiledMap.Layers[0].Tiles)

		fmt.Printf("Reading entities from layer \"%s\" \n", tiledMap.Layers[1].Name)

		if pos, ok := stage.GetEnterPosition(tiledMap.Layers[1].Tiles); ok {
			stage.EnterPosition = pos
		}

		game.InitEntities(world, stage, types.EntityCategoryObject, tiledMap.Layers[1].Tiles)
		game.InitEntities(world, stage, types.EntityCategoryEnemy, tiledMap.Layers[2].Tiles)

		if len(tiledMap.ObjectGroups) > 0 {
			game.InitObjectsEntities(world, stage, types.EntityCategoryObject, tiledMap.ObjectGroups[0])
		}

		outputDir := filepath.Join(StagesRootPath, stage.Name)

		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("error creating output dir: %v", err)
		}

		stage.Save(outputDir)
	}
}
