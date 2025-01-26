package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type actor interface {
	update(deltaTime float32)
	render()
}

type world struct {
	MapSize
	SoundPlayer
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
		//		rl.DrawTextureRec(w.blockTextures, rl.NewRectangle(float32(block.blockType)*w.blockWidth, 0, w.blockWidth, w.blockHeight), rl.NewVector2(float32(block.position.X)*w.blockWidth, float32(block.position.Y)*w.blockHeight), rl.White)
		block.render()
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
				po.Pushed(player, pushablePosition.Add(BlockPosition{X: offset.X, Y: 0}))

				//player.PushTo(po, pushablePosition.Add(BlockPosition{X: offset.X, Y: 0}))
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
				w.CheckForCollision(bo)

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

func (w *world) CheckForCollision(bo GroundObject) {

	collider, ok := bo.(Collider)

	if !ok {
		return
	}

	if collider.Body().IsColliding(w.player.Body()) {
		fmt.Println("Object colliding with Player")
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
