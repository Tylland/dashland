package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/systems"
)

const (
	StageRootPath = `..\..\assets\stages`
)

const assetsBase = `..\..\assets`

type DashlandGame struct {
	game.Screen
	Sounds
	camera game.Camera
	world  *ecs.World
	stage  *game.Stage
	player *ecs.Entity
}

func NewGame(screenWidth int, screenHeight int) *DashlandGame {
	game := DashlandGame{Screen: game.Screen{Width: screenWidth, Height: screenHeight}, world: ecs.NewWorld()}

	game.LoadSounds(filepath.Join(assetsBase, "sounds"))

	return &game
}

func (g *DashlandGame) Unload() {
	g.UnloadSounds()
	assets.UnloadAll()
}

func (g *DashlandGame) mapPath() string {
	return filepath.Join(assetsBase, "maps")
}

func (g *DashlandGame) mapFile(name string) string {
	return filepath.Join(g.mapPath(), name+".tmx")
}

func (g *DashlandGame) LoadTilesetTexture(tiledMap *tiled.Map, name string) (*rl.Texture2D, error) {
	for _, tileset := range tiledMap.Tilesets {
		if tileset.Name == name {
			path := path.Join("maps", tileset.Image.Source)
			name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			return assets.LoadTexture(name), nil
		}
	}

	return nil, fmt.Errorf("Layer %s not found!", name)
}

// func (g *DashlandGame) loadStageFromFile(name string) (*game.Stage, error) {
// 	tiledMap, err := tiled.LoadFile(g.mapFile(name))

// 	if err != nil {
// 		return nil, err
// 	}

// 	blockTexture, _ := g.LoadTilesetTexture(tiledMap, "Blocks")
// 	entityTextures, _ := g.LoadTilesetTexture(tiledMap, "Entities")

// 	groundCorners := assets.LoadTexture("ground_corners")

// 	mapSize := game.NewMapSize(tiledMap.Width, tiledMap.Height, tiledMap.TileWidth, tiledMap.TileHeight)
// 	stage := game.NewStage(name, mapSize, blockTexture, entityTextures, groundCorners)
// 	//	stage := &game.Stage{MapSize: mapSize, SoundPlayer: &g.Sounds, BlockMap: game.NewBlockMap(mapSize, blockTexture), EntityMap: game.NewEntityMap(mapSize, entityTextures, groundCorners)}

// 	if tiledMap.Properties != nil {
// 		stage.DiamondsRequired = tiledMap.Properties.GetInt("DiamondsRequired")
// 		stage.DiamondPoints = tiledMap.Properties.GetInt("DiamondPoints")
// 		stage.DiamondBonusPoints = tiledMap.Properties.GetInt("DiamondBonusPoints")
// 		stage.SetTimeLimit(tiledMap.Properties.GetInt("TimeLimit"))
// 	}

// 	fmt.Printf("Reading blocks from layer %s \n", tiledMap.Layers[0].Name)
// 	game.InitBlocks(g.world, stage, tiledMap.Layers[0].Tiles)

// 	fmt.Printf("Reading entities from layer \"%s\" \n", tiledMap.Layers[1].Name)

// 	if pos, ok := stage.GetEnterPosition(tiledMap.Layers[1].Tiles); ok {
// 		stage.EnterPosition = pos
// 	}

// 	game.InitEntities(g.world, stage, game.EntityCategoryObject, tiledMap.Layers[1].Tiles)
// 	game.InitEntities(g.world, stage, game.EntityCategoryEnemy, tiledMap.Layers[2].Tiles)

// 	if len(tiledMap.ObjectGroups) > 0 {
// 		game.InitObjectsEntities(g.world, stage, game.EntityCategoryObject, tiledMap.ObjectGroups[0])
// 	}

// 	return stage, nil
// }

// func (g *DashlandGame) LoadStage(name string, position common.BlockPosition) error {
// 	g.world.Clear()

// 	stage, err := g.loadStageFromFile(name)

// 	if err != nil {
// 		return err
// 	}

// 	if position.IsZero() {
// 		position = stage.EnterPosition
// 	}

// 	g.setupStage(stage, position)

// 	return nil
// }

func (g *DashlandGame) setupStage(stage *game.Stage, position common.BlockPosition) {
	playerTexture := assets.LoadTexture("player")

	player := game.NewPlayerEntity(g.world, stage, position, playerTexture)
	g.world.AddEntityNamed("player", player)

	playerPosition := ecs.GetComponent[components.PositionComponent](player)

	g.camera = game.NewSmoothFollowCamera(&g.Screen, playerPosition)

	g.world.AddSystems(
		systems.NewPauseSystem(),
		systems.NewStageTimerSystem(stage, &g.Sounds),
		systems.NewInputSystem(),
		systems.NewInputBehavior(stage),
		systems.NewGravityBehavior(stage),
		systems.NewWallWalkerBehavior(stage),
		systems.NewPushBehavior(stage),
		systems.NewBlockCollisionSystem(stage),
		systems.NewCollectorSystem(stage, &g.Sounds),
		systems.NewGameplaySystem(stage, &g.Sounds),
		systems.NewBlockMovement(stage),
		systems.NewAnimationSystem(),
		systems.NewParticleSystem(),
		systems.NewLifecycleSystem(),
	)

	g.world.AddSystem(systems.NewRenderSystem(stage, g.camera))
	g.world.AddSystem(systems.NewCleanup(stage, g))

	g.stage = stage
}

func (g *DashlandGame) Update(deltaTime float32) {
	g.world.Update(deltaTime)
	g.camera.Update(deltaTime)
}

func (g *DashlandGame) saveStage(stage *game.Stage) {

	outputDir := filepath.Join(`c:\\temp`, "dashland", "stages", stage.Name)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("error creating output dir: %v", err)
	}

	stage.Save(outputDir)
}

func (g *DashlandGame) LoadStage(name string, position common.BlockPosition) error {
	g.world.Clear()

	path := filepath.Join(StageRootPath, name)

	stage, err := g.loadStage(path)

	if err != nil {
		return err
	}

	if position.IsZero() {
		position = stage.EnterPosition
	}

	g.setupStage(stage, position)

	return nil
}

func (g *DashlandGame) loadStage(path string) (*game.Stage, error) {

	blockTexture := assets.LoadTexture("blocks")
	entityTextures := assets.LoadTexture("entities")
	groundCorners := assets.LoadTexture("ground_corners")

	mapSize := common.MapSize{}
	stage := game.NewStage("generated", mapSize, blockTexture, entityTextures, groundCorners)

	err := stage.Load(g.world, path)

	if err != nil {
		return nil, err
	}

	return stage, nil
}
