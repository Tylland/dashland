package game

import (
	"fmt"
	"path"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/game/core"
)

const mapPath = "maps/start.tmx" // Path to your Tiled Map.

type DashlandGame struct {
	Screen
	Sounds
	Camera Camera
	player *Player
	//	BlockMap BlockMap
	world *world
}

func NewGame(screenWidth int, screenHeight int) *DashlandGame {
	game := DashlandGame{Screen: Screen{Width: screenWidth, Height: screenHeight}}
	game.LoadSounds("sounds/effects")
	game.init()

	return &game
}

func (g *DashlandGame) Unload() {
	g.UnloadSounds()
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

func (g *DashlandGame) LoadTextureFromFile(source string) rl.Texture2D {
	fileName := path.Join(path.Dir(mapPath), source)

	return rl.LoadTexture(fileName)
}

func (g *DashlandGame) LoadWorldFromFile(filepath string) (*world, error) {

	tiledMap, err := tiled.LoadFile(filepath)

	if err != nil {
		return nil, err
	}

	blockTexture := g.LoadTextureFromFile(tiledMap.Tilesets[0].Image.Source)
	groundCorners := g.LoadTextureFromFile(tiledMap.Tilesets[1].Image.Source)

	fmt.Print(blockTexture)

	mapSize := MapSize{width: tiledMap.Width, height: tiledMap.Height, blockWidth: float32(tiledMap.TileWidth), blockHeight: float32(tiledMap.TileHeight)}
	world := &world{MapSize: mapSize, SoundPlayer: &g.Sounds, BlockMap: &BlockMap{MapSize: mapSize}, GroundMap: &GroundMap{MapSize: mapSize, entities: []*Entity{}}}
	world.RenderSystem = NewRenderSystem(world)

	world.blockTextures = blockTexture
	world.groundCorners = groundCorners

	fmt.Printf("Reading blocks from layer %s \n", tiledMap.Layers[0].Name)
	world.InitBlocks(world, tiledMap.Layers[0].Tiles)

	world.objectTextures = blockTexture

	fmt.Printf("Reading entities from layer \"%s\" \n", tiledMap.Layers[1].Name)
	//world.InitObjects(world, tiledMap.Layers[1].Tiles)

	//world.InitEntities(world, tiledMap.Layers[0].Tiles)
	world.InitEntities(world, tiledMap.Layers[1].Tiles)

	return world, nil
}

func (g *DashlandGame) init() {

	world, err := g.LoadWorldFromFile(mapPath)

	if err != nil {
		return
	}

	g.world = world

	g.player = NewPlayer(g)
	g.player.InitPosition(core.BlockPosition{X: 27, Y: 2})

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
