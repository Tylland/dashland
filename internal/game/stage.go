package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type SoundPlayer interface {
	PlayFx(name string)
}

type Stage struct {
	MapSize
	SoundPlayer
	*BlockMap
	*GroundMap
	player *Player
	//	actors []actor
}

func New() *Stage {
	w := &Stage{GroundMap: &GroundMap{entities: []*ecs.Entity{}}}

	return w
}

func (s *Stage) InitPlayer(player *Player) {
	s.player = player
	//	w.addActor(player)
}

func (s *Stage) Player() *Player {
	return s.player
}

func (s *Stage) GetPosition(position common.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*s.BlockMap.BlockWidth, float32(position.Y)*s.BlockMap.BlockHeight)
}

// func (w *Stage) update(deltaTime float32) {
// 	for _, block := range w.blocks {
// 		block.update(deltaTime)
// 	}

// 	// for _, obj := range w.objects {
// 	// 	if obj != nil {
// 	// 		obj.update(deltaTime)
// 	// 	}
// 	// }
// 	for _, act := range w.actors {
// 		act.update(deltaTime)
// 	}
// }

func (w *Stage) Render(deltaTime float32) {
	for _, block := range w.blocks {
		block.Render()
	}

}

// func (w *Stage) addActor(actor actor) {
// 	w.actors = append(w.actors, actor)
// }

func (s *Stage) IsObstacleForPlayer(world *ecs.World, player *Player, position common.BlockPosition) bool {
	block, success := s.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	if block.HasCharacteristic(characteristics.PlayerObstacle) {
		return true
	}

	entity := s.GetEntity(position)

	if entity == nil {
		return false
	}

	comps := world.GetComponents(entity)
	collectable := ecs.GetComponent[components.CollectableComponent](comps)

	if collectable != nil {
		return false
	}

	characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)

	if characteristic != nil {
		return false
	}

	if characteristic.Has(characteristics.Pushable) {
		// Calculate push direction based on player's position
		pushPos := position
		if player.Position.PreviousBlockPosition.X > position.X {
			pushPos = pushPos.Offset(-1, 0)
		} else if player.Position.PreviousBlockPosition.X < position.X {
			pushPos = pushPos.Offset(1, 0)
		}

		// Check if push position is free
		if s.CheckBlockAtPosition(Void, pushPos) && s.GetEntity(pushPos) == nil {
			return false
		}
	}

	return true
}

func (s *Stage) VisitBlock(position common.BlockPosition) {
	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, s.blocks[position.Y*s.Width+position.X].BlockType)
	s.SetBlock(NewBlock(s.BlockMap, s.GroundMap, Void, position), position)
	fmt.Printf(" to %d \n", s.blocks[position.Y*s.Width+position.X].BlockType)
}

func (s *Stage) CheckCharacteristics(world *ecs.World, position common.BlockPosition, character characteristics.Characteristics) bool {

	block, ok := s.GetBlock(position.X, position.Y)

	if !ok {
		return false
	}

	if block.HasCharacteristic(character) {
		return true
	}

	entity := s.GetEntity(position)

	if entity == nil {
		return false
	}

	characteristics := ecs.GetComponent[components.CharacteristicComponent](world.GetComponents(entity))

	return characteristics != nil && characteristics.Has(character)
}

func (s *Stage) CheckPositionOccupied(position common.BlockPosition) bool {
	if !s.CheckBlockAtPosition(Void, position) {
		return true
	}

	return s.GetEntity(position) != nil
}

func (s *Stage) CheckPlayerAtPosition(position common.BlockPosition) bool {
	if s.player.IsDead {
		return false
	}

	return s.player.Position.CurrentBlockPosition.IsSame(position)
}

func (s *Stage) OnEvent(event GameEvent) {
	switch e := event.(type) {
	case EntityCollisionEvent:
		fmt.Println("Entity collision detected!!")

		if e.Entity1.Type == EntityBoulder && e.Entity2.Type == EntityEnemy {
			s.OnBoulderEnemyCollision(event.World(), e.Entity1, e.Entity2)
		}
		if e.Entity2.Type == EntityBoulder && e.Entity1.Type == EntityEnemy {
			s.OnBoulderEnemyCollision(event.World(), e.Entity2, e.Entity1)
		}

	case BlockCollisionEvent:
		fmt.Println("Block collision detected!!")

		if e.Entity.Type == EntityBoulder && e.Block.BlockType == Soil {
			// Boulder is falling on player
			s.SoundPlayer.PlayFx("player_hurt")
		}

	case PlayerCollisionEvent:
		if e.Player.IsDead {
			return
		}

		fmt.Println("Player collision detected!!")

		if e.Entity.Type == EntityBoulder && e.EntityFalling {
			fmt.Println("Player hit by boulder!!")
			//Boulder is falling on player
			s.SoundPlayer.PlayFx("player_hurt")
			s.player.Hurt(e.Entity)
		}

		comps := e.World().GetComponents(e.Entity)
		character := ecs.GetComponent[components.CharacteristicComponent](comps)

		if character != nil && character.Has(characteristics.IsEnemy) {
			fmt.Println("Player hit by enemy!!")

			s.SoundPlayer.PlayFx("player_hurt")
			s.player.Hurt(e.Entity)
		}

		if e.Entity.Type == EntityDiamond {
			if e.EntityFalling {
				fmt.Println("Player hit by falling diamond!!")
				s.SoundPlayer.PlayFx("player_hurt")
				s.player.Hurt(e.Entity)
			} else {
				fmt.Println("Player collected diamond!!")
				s.SoundPlayer.PlayFx("diamond_collected")
				position := ecs.GetComponent[components.PositionComponent](comps)
				if position != nil {
					s.RemoveEntity(e.Entity, position)
				}
			}
		}

	}
}

func (s *Stage) OnBoulderEnemyCollision(world *ecs.World, boulder *ecs.Entity, enemy *ecs.Entity) {
	fmt.Println("Boulder and enemy collision detected!!")

	boulderComps := world.GetComponents(boulder)
	boulderVelocity := ecs.GetComponent[components.VelocityComponent](boulderComps)

	if boulderVelocity.BlockVector.IsZero() {
		return
	}

	enemyPosition := ecs.GetComponent[components.PositionComponent](world.GetComponents(enemy))

	s.RemoveEntity(enemy, enemyPosition)

	boulderPosition := ecs.GetComponent[components.PositionComponent](boulderComps)

	s.RemoveEntity(boulder, boulderPosition)

	position := enemyPosition.CurrentBlockPosition

	s.CreateDiamonds(world, position, 2, 2)
}

func (s *Stage) CreateDiamonds(world *ecs.World, position common.BlockPosition, width int, height int) {
	for y := -height; y <= height; y++ {
		for x := -width; x <= width; x++ {
			diamondPosition := position.Offset(x, y)
			if !s.CheckBlockAtPosition(Bedrock, diamondPosition) {
				s.SetBlock(NewBlock(s.BlockMap, s.GroundMap, Void, diamondPosition), diamondPosition)
				entity := NewDiamond(world, s, diamondPosition, s.GetPosition(diamondPosition))
				s.SetEntity(entity, diamondPosition)
			}
		}
	}
}

func (s *Stage) InitEntities(world *ecs.World, category ecs.EntityCategory, tiles []*tiled.LayerTile) {
	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.GroundMap.Width)

		entity, err := NewGameEntity(world, s, ecs.EntityType(uint32(category)+tile.ID), blockPosition)

		if err == nil {
			s.GroundMap.SetEntity(entity, blockPosition)
		}
	}
}

func (s *Stage) InitBlocks(world *Stage, tiles []*tiled.LayerTile) {
	s.blocks = make([]*Block, len(tiles))

	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.BlockMap.Width)
		s.BlockMap.SetBlock(NewBlock(s.BlockMap, s.GroundMap, BlockType(tile.ID), blockPosition), blockPosition)
	}

	s.BlockMap.PrintBlockMap()
}
