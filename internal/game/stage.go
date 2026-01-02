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
	*EntityMap
}

func New(enemyTextures *rl.Texture2D) *Stage {
	w := &Stage{EntityMap: &EntityMap{entities: []*ecs.Entity{}}}

	return w
}

func (s *Stage) GetPosition(position common.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*s.BlockMap.BlockWidth, float32(position.Y)*s.BlockMap.BlockHeight)
}

func (s *Stage) GetBlockPosition(position rl.Vector2) common.BlockPosition {
	bx := int(position.X / float32(s.BlockWidth))
	by := int(position.Y / float32(s.BlockHeight))

	return common.BlockPosition{X: bx, Y: by}
}

func (w *Stage) Render(deltaTime float32) {
	for _, block := range w.blocks {
		block.Render()
	}

}

func (s *Stage) VisitBlock(position common.BlockPosition) {
	fmt.Printf("Block at position %d,%d changed from %s", position.X, position.Y, s.blocks[position.Y*s.Width+position.X].BlockType.String())
	s.SetBlock(NewBlock(s.BlockMap, s.EntityMap, Void, position), position)
	fmt.Printf(" to %s \n", s.blocks[position.Y*s.Width+position.X].BlockType.String())
}

func (s *Stage) CheckCharacteristics(position common.BlockPosition, character characteristics.Characteristics) bool {

	if block, ok := s.GetBlock(position.X, position.Y); ok {
		if block.HasCharacteristic(character) {
			return true
		}
	}

	if entity, ok := s.GetEntityAtPosition(position); ok {
		characteristics := ecs.GetComponent[components.CharacteristicComponent](entity)

		return characteristics != nil && characteristics.Has(character)
	}

	return false
}

func (s *Stage) CheckBlocked(world *ecs.World, position common.BlockPosition, collider *components.ColliderComponent) bool {

	if block, ok := s.GetBlock(position.X, position.Y); ok {
		blocked, _ := collider.Result(block.Collider)
		if blocked {
			return true
		}
	}

	if entity, ok := s.GetEntityAtPosition(position); ok {
		entCharacter := ecs.GetComponent[components.CharacteristicComponent](entity)
		entCollider := ecs.GetComponent[components.ColliderComponent](entity)

		if entCharacter != nil && entCollider != nil {
			blocked, _ := collider.Result(entCollider)
			if blocked {
				return true
			}
		}
	}

	return false
}

func (s *Stage) CheckPositionOccupied(position common.BlockPosition) bool {
	if !s.CheckBlockAtPosition(Void, position) {
		return true
	}

	return s.GetEntity(position) != nil
}

// func (s *Stage) OnEvent(event any) {

// 	var world *ecs.World

// 	switch e := event.(type) {
// 	case *EntityCollisionEvent:
// 		fmt.Println("Entity collision detected!!")

// 		if e.Entity1.Type == EntityBoulder && e.Entity2.Type == EntityEnemy {
// 			s.OnBoulderEnemyCollision(world, e.Entity1, e.Entity2)
// 		}
// 		if e.Entity2.Type == EntityBoulder && e.Entity1.Type == EntityEnemy {
// 			s.OnBoulderEnemyCollision(world, e.Entity2, e.Entity1)
// 		}

// 	case *BlockCollisionEvent:
// 		fmt.Println("Block collision detected!!")

// 		if e.Entity.Type == EntityBoulder && e.Block.BlockType == Soil {
// 			// Boulder is falling on player
// 			s.SoundPlayer.PlayFx("player_hurt")
// 		}

// 	case *PlayerCollisionEvent:
// 		fmt.Println("Player collision detected!!")

// 		if e.Entity.Type == EntityBoulder && e.EntityFalling {
// 			fmt.Println("Player hit by boulder!!")
// 			//Boulder is falling on player
// 			s.SoundPlayer.PlayFx("player_hurt")
// 			//			s.player.Hurt(e.Entity)
// 		}

// 		character := ecs.GetComponent[components.CharacteristicComponent](e.Entity)

// 		if character != nil && character.Has(characteristics.IsEnemy) {
// 			fmt.Println("Player hit by enemy!!")

// 			s.SoundPlayer.PlayFx("player_hurt")
// 			//			s.player.Hurt(e.Entity)
// 		}

// 		if e.Entity.Type == EntityDiamond {
// 			if e.EntityFalling {
// 				fmt.Println("Player hit by falling diamond!!")
// 				s.SoundPlayer.PlayFx("player_hurt")
// 				//				s.player.Hurt(e.Entity)
// 			} else {
// 				fmt.Println("Player collected diamond!!")
// 				s.SoundPlayer.PlayFx("diamond_collected")
// 				position := ecs.GetComponent[components.PositionComponent](e.Entity)
// 				if position != nil {
// 					s.RemoveEntity(e.Entity, position.CurrentBlockPosition)
// 				}
// 			}
// 		}
// 	default:
// 		fmt.Printf("Unknown event type %T \n", e)
// 	}
// }

// func (s *Stage) OnBoulderEnemyCollision(world *ecs.World, boulder *ecs.Entity, enemy *ecs.Entity) {
// 	fmt.Println("Boulder and enemy collision detected!!")

// 	boulderVelocity := ecs.GetComponent[components.VelocityComponent](boulder)

// 	if !boulderVelocity.IsMoving() {
// 		return
// 	}

// 	enemyPosition := ecs.GetComponent[components.PositionComponent](enemy)

// 	s.RemoveEntity(enemy, enemyPosition.CurrentBlockPosition)

// 	boulderPosition := ecs.GetComponent[components.PositionComponent](boulder)

// 	s.RemoveEntity(boulder, boulderPosition.CurrentBlockPosition)

// 	position := enemyPosition.CurrentBlockPosition

// 	s.CreateDiamonds(world, position, 2, 2)
// }

// func (s *Stage) CreateDiamonds(world *ecs.World, position common.BlockPosition, width int, height int) {
// 	for y := -height; y <= height; y++ {
// 		for x := -width; x <= width; x++ {
// 			diamondPosition := position.Offset(x, y)
// 			if !s.CheckBlockAtPosition(Bedrock, diamondPosition) {
// 				s.SetBlock(NewBlock(s.BlockMap, s.EntityMap, Void, diamondPosition), diamondPosition)
// 				entity := NewDiamond(world, s, diamondPosition, s.GetPosition(diamondPosition))
// 				s.SetEntity(entity, diamondPosition)
// 			}
// 		}
// 	}
// }

func (s *Stage) InitEntities(world *ecs.World, category ecs.EntityCategory, tiles []*tiled.LayerTile) {
	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.EntityMap.Width)

		entity, err := NewGameEntity(world, s, ecs.EntityType(uint32(category)+tile.ID), blockPosition)

		if err == nil {
			s.EntityMap.SetEntity(entity, blockPosition)
		}
	}
}

func (s *Stage) InitBlocks(world *Stage, tiles []*tiled.LayerTile) {
	s.blocks = make([]*Block, len(tiles))

	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.BlockMap.Width)
		s.BlockMap.SetBlock(NewBlock(s.BlockMap, s.EntityMap, BlockType(tile.ID), blockPosition), blockPosition)
	}

	s.BlockMap.PrintBlockMap()
}

func getString(obj *tiled.Object, name string) string {
	value := ""
	for _, p := range obj.Properties {
		if p.Name == name {
			var val any = p.Value
			switch v := val.(type) {
			case string:
				value = v
			case float64:
				value = fmt.Sprintf("%v", v)
			default:
				value = fmt.Sprintf("%v", v)
			}
			break
		}
	}

	return value
}

func (s *Stage) InitObjectsEntities(world *ecs.World, category ecs.EntityCategory, objectLayer *tiled.ObjectGroup) {
	for _, obj := range objectLayer.Objects {
		if obj.Type == "EntityDoor" {

			blockPos := s.GetBlockPosition(rl.Vector2{X: float32(obj.X), Y: float32(obj.Y)})
			targetStage := getString(obj, "Stage")
			targetPosition, err := common.ParseBlockPosition(getString(obj, "Position"))

			if err != nil {
				panic(err.Error())
			}

			door := NewDoorWithDestination(world, s, blockPos, s.GetPosition(blockPos), targetStage, targetPosition)
			s.EntityMap.SetEntity(door, blockPos)
		}
	}
}
