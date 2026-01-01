package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type GameplaySystem struct {
	stage *game.Stage
	sound game.SoundPlayer
}

func NewGameplaySystem(stage *game.Stage, sound game.SoundPlayer) *GameplaySystem {
	return &GameplaySystem{stage: stage, sound: sound}
}

func (s *GameplaySystem) Update(world *ecs.World, deltaTime float32) {
	for _, event := range world.Events() {
		if event != nil && event.Name == "damage" {
			s.handleDamage(world, event.Data.(*game.DamageEvent))
		}

		if event != nil && event.Name == "blockcollision" {
			s.handleBlockCollision(world, event.Data.(*game.BlockCollisionEvent))
		}

		if event != nil && event.Name == "entitycollision" {
			s.handleCollision(world, event.Data.(*game.EntityCollisionEvent))
		}
	}
}

func (s *GameplaySystem) handleBlockCollision(world *ecs.World, collision *game.BlockCollisionEvent) {
	if collision.Block.BlockType == game.Soil && collision.Entity.Type == game.EntityPlayer {
		poistion := ecs.GetComponent[components.PositionComponent](collision.Entity)
		step := ecs.GetComponent[components.BlockStep](collision.Entity)

		s.stage.VisitBlock(poistion.CurrentBlockPosition.Add(step.Increment))
	}

	if collision.Block.BlockType == game.Soil && collision.Entity.Type == game.EntityBoulder {
		s.sound.PlayFx("boulder_collision")
	}

	if collision.Block.BlockType == game.Soil && collision.Entity.Type == game.EntityDiamond {
		s.sound.PlayFx("diamond_collision")
	}
}

func (s *GameplaySystem) handleCollision(world *ecs.World, collision *game.EntityCollisionEvent) {
	// if collision.Entity1.Type == game.EntityPlayer && collision.Entity2.Type == game.EntityEnemy {
	// 	s.OnBoulderEnemyCollision(world, collision.Entity1, collision.Entity2)
	// }

	// if collision.Entity1.Type == game.EntityBoulder && collision.Entity2.Type == game.EntityEnemy {
	// 	s.OnBoulderEnemyCollision(world, collision.Entity1, collision.Entity2)
	// }

	// if collision.Entity2.Type == game.EntityBoulder && collision.Entity1.Type == game.EntityEnemy {
	// 	s.OnBoulderEnemyCollision(world, collision.Entity2, collision.Entity1)
	// }
}

func (s *GameplaySystem) handleDamage(world *ecs.World, damage *game.DamageEvent) {

	if damage.Target.Type == game.EntityPlayer {
		s.DamageOnPlayer(world, damage.Source, damage.Target)
	}

	if damage.Target.Type == game.EntityEnemy {
		s.OnDamageOnEnemy(world, damage.Source, damage.Target)
	}
}

func (s *GameplaySystem) OnDamageOnEnemy(world *ecs.World, boulder *ecs.Entity, enemy *ecs.Entity) {
	fmt.Println("Boulder and enemy collision detected!!")

	boulderStep := ecs.GetComponent[components.BlockStep](boulder)

	if boulderStep.Direction.IsZero() {
		return
	}

	enemyPosition := ecs.GetComponent[components.PositionComponent](enemy)

	s.stage.TryRemoveEntity(enemy)
	world.EnqueueRemoval(enemy)

	s.stage.TryRemoveEntity(boulder)
	world.EnqueueRemoval(boulder)

	position := enemyPosition.CurrentBlockPosition

	s.CreateDiamonds(world, position, 2, 2)
}

func (s *GameplaySystem) OnBoulderPlayerCollision(world *ecs.World, boulder *ecs.Entity, enemy *ecs.Entity) {
	fmt.Println("Boulder and enemy collision detected!!")

	boulderStep := ecs.GetComponent[components.BlockStep](boulder)

	if boulderStep.Direction.IsZero() {
		return
	}

	enemyPosition := ecs.GetComponent[components.PositionComponent](enemy)

	s.stage.TryRemoveEntity(enemy)
	world.EnqueueRemoval(enemy)

	s.stage.TryRemoveEntity(boulder)
	world.EnqueueRemoval(boulder)

	position := enemyPosition.CurrentBlockPosition

	s.CreateDiamonds(world, position, 2, 2)
}

func (s *GameplaySystem) DamageOnPlayer(world *ecs.World, enemy *ecs.Entity, player *ecs.Entity) {
	fmt.Println("Player damage triggered!!")

	s.sound.PlayFx("player_hurt")

	playerPosition := ecs.GetComponent[components.PositionComponent](player)
	position := playerPosition.CurrentBlockPosition

	world.EnqueueRemoval(player)

	s.CreateDiamonds(world, position, 2, 2)
}

func (s *GameplaySystem) CreateDiamonds(world *ecs.World, position common.BlockPosition, width int, height int) {
	for y := -height; y <= height; y++ {
		for x := -width; x <= width; x++ {
			diamondPosition := position.Offset(x, y)
			if s.stage.CheckCharacteristics(diamondPosition, characteristics.Destructable) {
				s.stage.SetBlock(game.NewBlock(s.stage.BlockMap, s.stage.EntityMap, game.Void, diamondPosition), diamondPosition)

				doomed := s.stage.GetEntity(diamondPosition)
				if doomed != nil {
					world.EnqueueRemoval(doomed)
				}

				diamond, err := game.NewGameEntity(world, s.stage, game.EntityDiamond, diamondPosition)

				if err == nil {
					s.stage.SetEntity(diamond, diamondPosition)
				}
			}
		}
	}
}
