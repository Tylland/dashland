package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/core"
)

type actor interface {
	update(deltaTime float32)
	render()
}

type World struct {
	MapSize
	SoundPlayer
	*BlockMap
	*GroundMap
	player       *Player
	actors       []actor
	RenderSystem *RenderSystem
}

func NewWorld() *World {
	w := &World{actors: []actor{}, GroundMap: &GroundMap{entities: []*Entity{}}}

	w.RenderSystem = NewRenderSystem(w)

	return w
}

func (w *World) initPlayer(player *Player) {
	w.player = player
	w.addActor(player)
}

func (w *World) GetPosition(position core.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*w.blockWidth, float32(position.Y)*w.blockHeight)
}

func (w *World) update(deltaTime float32) {
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
	NewWallWalkerSystem(w).Update()
	NewBlockMovementSystem(w).Update(deltaTime)
	NewBlockCollisionSystem(w).Update()

}

func (w *World) render(deltaTime float32) {
	for _, block := range w.blocks {
		block.render()
	}

	w.RenderSystem.Update(deltaTime)

	for _, act := range w.actors {
		act.render()
	}
}

// Entities exposes the ground map entities for systems that iterate over them.
func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) addActor(actor actor) {
	w.actors = append(w.actors, actor)
}

func (w *World) IsObstacleForPlayer(player *Player, position core.BlockPosition) bool {
	block, success := w.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	if block.HasCharacteristic(characteristics.PlayerObstacle) {
		return true
	}

	entity := w.GetEntity(position)

	if entity == nil {
		return false
	}

	if entity.Collectable != nil {
		return false
	}

	if entity.HasCharacteristic(characteristics.Pushable) {
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

func (w *World) VisitBlock(position core.BlockPosition) {
	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, w.blocks[position.Y*w.width+position.X].blockType)
	w.SetBlock(NewBlock(w, Void, position.X, position.Y), position)
	fmt.Printf(" to %d \n", w.blocks[position.Y*w.width+position.X].blockType)
}

func (w *World) CheckCharacteristics(position core.BlockPosition, characteristics characteristics.Characteristics) bool {

	block, ok := w.GetBlock(position.X, position.Y)

	if !ok {
		return false
	}

	if block.HasCharacteristic(characteristics) {
		return true
	}

	entity := w.GetEntity(position)

	if entity == nil {
		return false
	}

	return entity.HasCharacteristic(characteristics)
}

func (w *World) checkPositionOccupied(position core.BlockPosition) bool {
	if !w.CheckBlockAtPosition(Void, position) {
		return true
	}

	return w.GetEntity(position) != nil
}

func (w *World) checkPlayerAtPosition(position core.BlockPosition) bool {
	if w.player.IsDead {
		return false
	}

	return w.player.Position.CurrentBlockPosition.IsSame(position)
}

func (w *World) GetNearbyBlocks(position *core.BlockPosition) []*Block {
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

func (w *World) OnEvent(event GameEvent) {
	switch e := event.(type) {
	case EntityCollisionEvent:
		fmt.Println("Entity collision detected!!")

		if e.Entity1.Type == EntityBoulder && e.Entity2.Type == EntityEnemy {
			w.OnBoulderEnemyCollision(e.Entity1, e.Entity2)
		}
		if e.Entity2.Type == EntityBoulder && e.Entity1.Type == EntityEnemy {
			w.OnBoulderEnemyCollision(e.Entity2, e.Entity1)
		}

	case BlockCollisionEvent:
		fmt.Println("Block collision detected!!")

		if e.Entity.Type == EntityBoulder && e.Block.blockType == Soil {
			// Boulder is falling on player
			w.SoundPlayer.PlayFx("player_hurt")
		}

	case PlayerCollisionEvent:
		if e.Player.IsDead {
			return
		}

		fmt.Println("Player collision detected!!")

		if e.Entity.Type == EntityBoulder && e.EntityFalling {
			fmt.Println("Player hit by boulder!!")
			// Boulder is falling on player
			w.SoundPlayer.PlayFx("player_hurt")
			w.player.Hurt(e.Entity)
		}

		if e.Entity.HasCharacteristic(characteristics.IsEnemy) {
			fmt.Println("Player hit by enemy!!")

			w.SoundPlayer.PlayFx("player_hurt")
			w.player.Hurt(e.Entity)
		}

		if e.Entity.Type == EntityDiamond {
			if e.EntityFalling {
				fmt.Println("Player hit by falling diamond!!")
				w.SoundPlayer.PlayFx("player_hurt")
				w.player.Hurt(e.Entity)
			} else {
				fmt.Println("Player collected diamond!!")
				w.SoundPlayer.PlayFx("diamond_collected")
				w.RemoveEntity(e.Entity)
			}
		}

	}
}

func (w *World) OnBoulderEnemyCollision(boulder *Entity, enemy *Entity) {
	fmt.Println("Boulder and enemy collision detected!!")

	if boulder.Velocity.BlockVector.IsZero() {
		return
	}

	w.RemoveEntity(enemy)
	w.RemoveEntity(boulder)

	position := enemy.Position.CurrentBlockPosition

	w.CreateDiamonds(position, 2, 2)
}

func (w *World) CreateDiamonds(position core.BlockPosition, width int, height int) {
	for y := -height; y <= height; y++ {
		for x := -width; x <= width; x++ {
			diamondPosition := position.Offset(x, y)
			if !w.CheckBlockAtPosition(Bedrock, diamondPosition) {
				w.SetBlock(NewBlock(w, Void, diamondPosition.X, diamondPosition.Y), diamondPosition)
				entity := NewDiamond(w, NewEntityId(), diamondPosition, w.GetPosition(diamondPosition))
				w.SetEntity(entity, diamondPosition)
			}
		}
	}
}
