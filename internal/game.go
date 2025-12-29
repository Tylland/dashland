package internal

import (
	"fmt"
	"path"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/systems"
)

const assetsBase = "assets"

type DashlandGame struct {
	game.Screen
	Sounds
	camera game.Camera
	world  *ecs.World
	stage  *game.Stage
}

func NewGame(screenWidth int, screenHeight int) *DashlandGame {
	game := DashlandGame{Screen: game.Screen{Width: screenWidth, Height: screenHeight}, world: ecs.NewWorld()}
	game.LoadSounds(filepath.Join(assetsBase, "sounds"))

	return &game
}

func (g *DashlandGame) Unload() {
	g.UnloadSounds()
}

func (g *DashlandGame) mapPath() string {
	return filepath.Join(assetsBase, "maps")
}

func (g *DashlandGame) mapFile(name string) string {
	return filepath.Join(g.mapPath(), name+".tmx")
}

// func (g *DashlandGame) CteateTextureFromFile(source string) *rl.Texture2D {
// 	fileName := path.Join(path.Dir(mapPath), source)

// 	file, err := os.Open(fileName)
// 	defer file.Close()

// 	if err != nil {
// 		return nil
// 	}

// 	buf := new(bytes.Buffer)

// 	var new_image image.Image

// 	png.Encode(buf, new_image)

// 	rl.NewImage(buf.Bytes())

// 	rl.LoadImage(fileName)

// 	return rl.LoadTexture(fileName)
// }

func (g *DashlandGame) LoadTextureFromAssets(source string) *rl.Texture2D {
	fileName := path.Join(assetsBase, "images", source)

	texture := rl.LoadTexture(fileName)
	return &texture
}

func (g *DashlandGame) LoadMapTexture(source string) *rl.Texture2D {
	fileName := filepath.Join(g.mapPath(), source)

	texture := rl.LoadTexture(fileName)
	return &texture
}

func (g *DashlandGame) LoadTilesetTexture(tiledMap *tiled.Map, name string) *rl.Texture2D {
	for _, tileset := range tiledMap.Tilesets {
		if tileset.Name == name {
			return g.LoadMapTexture(tileset.Image.Source)
		}
	}

	return nil
}

func (g *DashlandGame) loadStageFromFile(name string) (*game.Stage, error) {
	tiledMap, err := tiled.LoadFile(g.mapFile(name))

	if err != nil {
		return nil, err
	}

	blockTexture := g.LoadTilesetTexture(tiledMap, "Blocks")
	entityTextures := g.LoadTilesetTexture(tiledMap, "Entities")
	groundCorners := g.LoadTilesetTexture(tiledMap, "GroundCorners")

	enemyTexture := g.LoadTextureFromAssets("animations.png")

	mapSize := game.MapSize{Width: tiledMap.Width, Height: tiledMap.Height, BlockWidth: float32(tiledMap.TileWidth), BlockHeight: float32(tiledMap.TileHeight)}
	stage := &game.Stage{MapSize: mapSize, SoundPlayer: &g.Sounds, BlockMap: game.NewBlockMap(mapSize, blockTexture), EntityMap: game.NewGroundMap(mapSize, entityTextures, enemyTexture, groundCorners)}

	fmt.Printf("Reading blocks from layer %s \n", tiledMap.Layers[0].Name)
	stage.InitBlocks(stage, tiledMap.Layers[0].Tiles)

	fmt.Printf("Reading entities from layer \"%s\" \n", tiledMap.Layers[1].Name)

	stage.InitPlayerPosition(tiledMap.Layers[1].Tiles)

	stage.InitEntities(g.world, game.EntityCategoryObject, tiledMap.Layers[1].Tiles)
	stage.InitEntities(g.world, game.EntityCategoryEnemy, tiledMap.Layers[2].Tiles)

	return stage, nil
}

func (g *DashlandGame) LoadStage(name string) error {
	stage, err := g.loadStageFromFile(name)

	if err != nil {
		return err
	}

	playerTexture := g.LoadTextureFromAssets("player.png")

	player := game.NewPlayerEntity(g.world, stage, stage.InitialPlayerPosition, playerTexture)
	g.world.AddEntityNamed("player", player)

	playerPosition := ecs.GetComponent[components.PositionComponent](player)

	g.camera = game.NewSmoothFollowCamera(&g.Screen, playerPosition)

	g.world.AddSystems(
		systems.NewInputSystem(),
		systems.NewInputBehavior(stage),
		systems.NewGravityBehavior(stage),
		systems.NewWallWalkerBehavior(stage),
		systems.NewPushBehavior(stage),
		systems.NewBlockCollisionSystem(stage),
		systems.NewCollect(stage, stage.SoundPlayer),
		systems.NewGameplaySystem(stage, stage.SoundPlayer),
		systems.NewBlockMovement(stage),
		systems.NewAnimationSystem(),
	)

	g.world.AddSystem(systems.NewRenderSystem(stage, g.camera))
	g.world.AddSystem(systems.NewCleanup(stage))

	g.stage = stage

	return nil
}

func (g *DashlandGame) Update(deltaTime float32) {
	//fmt.Println(deltaTime)
	g.world.Update(deltaTime)

	//	g.player.Update(g.world, deltaTime)

	// g.BlockMap.Update(deltaTime)
	// g.Player.update(deltaTime)

	g.camera.Update(deltaTime)
}
