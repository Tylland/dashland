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
	MapSize
	*BlockMap
	*GroundMap
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

func LoadWorldFromFile(filepath string) (*world, error) {

	tiledMap, err := tiled.LoadFile(filepath)

	if err != nil {
		return nil, err
	}

	relativeImagePath := tiledMap.Tilesets[0].Image.Source

	fileName := path.Join(path.Dir(mapPath), relativeImagePath)

	blockTexture := rl.LoadTexture(fileName)

	fmt.Print(blockTexture)

	mapSize := MapSize{width: tiledMap.Width, height: tiledMap.Height, blockWidth: float32(tiledMap.TileWidth), blockHeight: float32(tiledMap.TileHeight)}
	world := &world{MapSize: mapSize, BlockMap: &BlockMap{MapSize: mapSize}, GroundMap: &GroundMap{MapSize: mapSize}}

	world.blockTextures = blockTexture

	fmt.Printf("Reading blocks from layer %s \n", tiledMap.Layers[0].Name)
	world.InitBlocks(world, tiledMap.Layers[0].Tiles)

	world.objectTextures = blockTexture

	fmt.Printf("Reading objects from layer \"%s\" \n", tiledMap.Layers[1].Name)
	world.InitObjects(world, tiledMap.Layers[1].Tiles)

	return world, nil
}

func (w *world) GetPosition(position BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*w.blockWidth, float32(position.Y)*w.blockHeight)
}

func (w *world) update(deltaTime float32) {
	for _, block := range w.blocks {
		block.update(deltaTime)
	}

	for _, obj := range w.objects {
		if obj != nil {
			obj.update(deltaTime)
		}
	}

	for _, act := range w.actors {
		act.update(deltaTime)
	}
}

func (w *world) render() {
	for _, block := range w.blocks {
		rl.DrawTextureRec(w.blockTextures, rl.NewRectangle(float32(block.blockType)*w.blockWidth, 0, w.blockWidth, w.blockHeight), rl.NewVector2(float32(block.position.X)*w.blockWidth, float32(block.position.Y)*w.blockHeight), rl.White)
	}

	for _, obj := range w.objects {
		if obj != nil {
			obj.render()
		}
	}

	for _, act := range w.actors {
		act.render()
	}
}

func (w *world) addActor(actor actor) {
	w.actors = append(w.actors, actor)
}

func (w *world) IsObstacleForPlayer(player *Player, position BlockPosition) bool {
	block, success := w.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	if block.IsObstacleForPlayer(player) {
		return true
	}

	object := w.GetObject(position)

	if object == nil {
		return false
	}

	return object.IsObstacleForPlayer(player)
}

func (w *world) VisitBlock(position BlockPosition) {
	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, w.blocks[position.Y*w.width+position.X].blockType)
	w.SetBlock(NewBlock(w, Void, position.X, position.Y), position)
	fmt.Printf(" to %d \n", w.blocks[position.Y*w.width+position.X].blockType)
}

func (w *world) VisitObject(player *Player, position BlockPosition) {

	obj := w.GetObject(position)

	if obj != nil {
		co, ok := obj.(CollectableObject)

		if ok && co != nil {
			player.Collect(co)
		}

		po, ok := obj.(PushableObject)

		if ok && po != nil {

			pushablePosition := po.GetBlockPosition()

			offset := pushablePosition.Subtract(player.blockPosition)

			if offset.Y == 0 {
				player.PushTo(po, pushablePosition.Add(BlockPosition{X: offset.X, Y: 0}))
			}
		}
	}
}

func (w *world) checkPositionOccupied(position BlockPosition) bool {
	if !w.CheckBlockAtPosition(Void, position) {
		return true
	}

	if w.player.blockPosition.IsSame(position) {
		return true
	}

	return w.GetObject(position) != nil
}

func (w *world) checkPlayerAtPosition(position BlockPosition) bool {
	return w.player.blockPosition.IsSame(position)
}

func (w *world) ApplyGravity(bo FallingObject, deltaTime float32) {
	fmt.Println("Boulder Upadtes")

	bo.UpdateFalling(deltaTime)

	if bo.HasBehavior(CanFall) && !bo.IsFalling() {
		current := bo.GetBlockPosition()

		under := current.Offset(0, 1)

		if !w.CheckBlockAtPosition(Soil, under) && !w.checkPlayerAtPosition(under) {

			if !w.checkPositionOccupied(under) {
				bo.StartFalling(w.GetPosition(current), w.GetPosition(under))
				w.MoveObject(bo, under)
				return
			}

			right := current.Offset(1, 0)
			rightUnder := current.Offset(1, 1)

			if !w.checkPositionOccupied(right) && !w.checkPositionOccupied(rightUnder) {
				bo.StartFalling(w.GetPosition(current), w.GetPosition(rightUnder))
				w.MoveObject(bo, rightUnder)
				return
			}

			left := current.Offset(-1, 0)
			leftUnder := current.Offset(-1, 1)

			if !w.checkPositionOccupied(left) && !w.checkPositionOccupied(leftUnder) {
				bo.StartFalling(w.GetPosition(current), w.GetPosition(leftUnder))
				w.MoveObject(bo, leftUnder)
				return
			}
		}

	}
}

// func (w *world) ApplyPush(bo PushableObject, deltaTime float32) {
// 	fmt.Println("Boulder Upadtes")

// 	bo.UpdateFalling(deltaTime)

// 	if bo.HasBehavior(CanFall) && !bo.IsFalling() {
// 		current := bo.getBlockPosition()

// 		under := current.Offset(0, 1)

// 		if !w.CheckBlockAtPosition(Soil, under) && !w.checkPlayerAtPosition(under) {

// 			if !w.checkPositionOccupied(under) {
// 				bo.StartFalling(w.GetPosition(current), w.GetPosition(under))
// 				w.MoveObject(bo, under)
// 				return
// 			}

// 			right := current.Offset(1, 0)
// 			rightUnder := current.Offset(1, 1)

// 			if !w.checkPositionOccupied(right) && !w.checkPositionOccupied(rightUnder) {
// 				bo.StartFalling(w.GetPosition(current), w.GetPosition(rightUnder))
// 				w.MoveObject(bo, rightUnder)
// 				return
// 			}

// 			left := current.Offset(-1, 0)
// 			leftUnder := current.Offset(-1, 1)

// 			if !w.checkPositionOccupied(left) && !w.checkPositionOccupied(leftUnder) {
// 				bo.StartFalling(w.GetPosition(current), w.GetPosition(leftUnder))
// 				w.MoveObject(bo, leftUnder)
// 				return
// 			}
// 		}

// 	}
// }
