package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/core"
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
	player       *Player
	actors       []actor
	RenderSystem *RenderSystem
}

func NewWorld() *world {
	w := &world{actors: []actor{}, GroundMap: &GroundMap{entities: []*Entity{}}}

	w.RenderSystem = NewRenderSystem(w)

	return w
}

func (w *world) initPlayer(player *Player) {
	w.player = player
	w.addActor(player)
}

func (w *world) GetPosition(position core.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*w.blockWidth, float32(position.Y)*w.blockHeight)
}

func (w *world) update(deltaTime float32) {
	for _, block := range w.blocks {
		block.update(deltaTime)
	}

	// for _, obj := range w.objects {
	// 	if obj != nil {
	// 		obj.update(deltaTime)
	// 	}
	// }
	for _, act := range w.actors {
		act.update(deltaTime)
	}

	NewGravitySystem(w).Update()
	NewBlockMovementSystem(w).Update(deltaTime)
	NewBlockCollisionSystem(w).Update()

}

func (w *world) render() {
	for _, block := range w.blocks {
		block.render()
	}

	w.RenderSystem.Update()

	for _, act := range w.actors {
		act.render()
	}
}

func (w *world) addActor(actor actor) {
	w.actors = append(w.actors, actor)
}

func (w *world) IsObstacleForPlayer(player *Player, position core.BlockPosition) bool {
	block, success := w.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	if block.IsObstacleForPlayer(player) {
		return true
	}

	entity := w.GetEntity(position)

	if entity == nil {
		return false
	}

	if entity.Collectable != nil {
		return false
	}

	if entity.Behavior&Pushable != 0 {
		// Calculate push direction based on player's position
		pushPos := position
		if player.Position.PreviousBlockPosition.X > position.X {
			pushPos = pushPos.Offset(-1, 0)
		} else if player.Position.PreviousBlockPosition.X < position.X {
			pushPos = pushPos.Offset(1, 0)
		}

		// Check if push position is free
		if w.CheckBlockAtPosition(Void, pushPos) && w.GetEntity(pushPos) == nil {
			return false
		}
	}

	return true
}

func (w *world) VisitBlock(position core.BlockPosition) {
	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, w.blocks[position.Y*w.width+position.X].blockType)
	w.SetBlock(NewBlock(w, Void, position.X, position.Y), position)
	fmt.Printf(" to %d \n", w.blocks[position.Y*w.width+position.X].blockType)
}

/*
func (w *world) VisitObject(player *Player, position core.BlockPosition) {

	obj := w.GetObject(position)

	if obj != nil {
		co, ok := obj.(CollectableObject)

		if ok && co != nil {
			player.Collect(co)
		}

		po, ok := obj.(PushableObject)

		if ok && po != nil {

			pushablePosition := po.GetBlockPosition()

			offset := pushablePosition.Subtract(player.Position.BlockPosition)

			if offset.Y == 0 {
				po.Pushed(player, pushablePosition.Add(core.BlockPosition{X: offset.X, Y: 0}))

				//player.PushTo(po, pushablePosition.Add(BlockPosition{X: offset.X, Y: 0}))
			}
		}
	}
}*/

func (w *world) checkPositionOccupied(position core.BlockPosition) bool {
	if !w.CheckBlockAtPosition(Void, position) {
		return true
	}

	// if w.player.Position.BlockPosition.IsSame(position) {
	// 	return true
	// }

	return w.GetEntity(position) != nil
}

func (w *world) checkPlayerAtPosition(position core.BlockPosition) bool {
	return w.player.Position.CurrentBlockPosition.IsSame(position)
}

/*func (w *world) ApplyGravity(bo FallingObject, deltaTime float32) {
	fmt.Println("Boulder Upadtes")

	bo.UpdateFalling(deltaTime)

	if bo.HasBehavior(CanFall) && !bo.Falling() {
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
}*/

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

func (w *world) GetNearbyBlocks(position *core.BlockPosition) []*Block {
	// Get entity's block position
	blockX := position.X
	blockY := position.Y

	// Check surrounding blocks (3x3 grid around entity)
	var nearbyBlocks []*Block

	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			checkX := blockX + x
			checkY := blockY + y

			// Check world bounds
			if checkX >= 0 && checkX < w.width &&
				checkY >= 0 && checkY < w.height {
				if block, ok := w.GetBlock(checkX, checkY); ok && block != nil {
					nearbyBlocks = append(nearbyBlocks, block)
				}
			}
		}
	}

	return nearbyBlocks
}

func (w *world) OnEvent(event GameEvent) {
	switch e := event.(type) {
	case EntityCollisionEvent:
		fmt.Println("Entity collision detected!!")

		if e.Entity1.Type == Boulder && e.Entity2.Type == Soil {
			e.Entity1.Velocity.Vector.Clear()
			e.Entity1.Position.Y = e.Entity2.Position.Y - e.Entity1.Collision.Height
		}
		if e.Entity1.Type == Soil && e.Entity2.Type == Boulder {
			e.Entity2.Velocity.Vector.Clear()
			e.Entity2.Position.Y = e.Entity1.Position.Y + e.Entity2.Collision.Height
		}

	case BlockCollisionEvent:
		fmt.Println("Block collision detected!!")

	case PlayerCollisionEvent:
		fmt.Println("Player collision detected!!")

		if e.Entity.Type == Boulder && e.EntityFalling {
			// Boulder is falling on player
			w.SoundPlayer.PlayFx("player_hurt")
			w.player.Hurt(e.Entity)
		}

		if e.Entity.Type == Diamond {
			if e.EntityFalling {
				w.SoundPlayer.PlayFx("player_hurt")
				w.player.Hurt(e.Entity)
			} else {
				w.SoundPlayer.PlayFx("diamond_collected")
				w.RemoveEntity(e.Entity)
			}
		}
	}
}
