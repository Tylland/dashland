package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
)

type GameplaySystem struct {
	stage *game.Stage
	sound game.SoundPlayer
}

func NewGameplaySystem(stage *game.Stage, sound game.SoundPlayer) *GameplaySystem {
	return &GameplaySystem{stage: stage, sound: sound}
}

func (s *GameplaySystem) Update(world *ecs.World, deltaTime float32) {
	magicWallWasActive := s.stage.IsMagicWallActive()
	s.stage.UpdateMagicWall(deltaTime)

	// Magic wall just expired → switch all magic wall entities to idle animation
	if magicWallWasActive && !s.stage.IsMagicWallActive() {
		s.setMagicWallAnimation(world, "idle")
	}

	for _, event := range world.Events() {
		if event != nil && event.Name == "exitopen" {
			s.handleExitOpen(world)
		}

		if event != nil && event.Name == "stageexpired" {
			s.stage.TimeExpired = true
		}

		if event != nil && event.Name == "timeoutdeath" {
			if player := world.GetEntity("player"); player != nil {
				s.DamageOnPlayer(world, player)
			}
		}

		if event != nil && event.Name == "damage" {
			s.handleDamage(world, event.Data.(*game.DamageEvent))
		}

		if event != nil && event.Name == "blockcollision" {
			s.handleBlockCollision(world, event.Data.(*game.BlockCollisionEvent))
		}

		if event != nil && event.Name == "entityobstructed" {
			s.handleEntityBlocks(world, event.Data.(*game.EntityEvent))
		}

		if event != nil && event.Name == "entitycollision" {
			s.handleCollision(world, event.Data.(*game.EntityCollisionEvent))
		}
	}
}

func (s *GameplaySystem) RunWhenPaused() bool {
	return true
}

func (s *GameplaySystem) handleBlockCollision(world *ecs.World, collision *game.BlockCollisionEvent) {
	if collision.Block.Type == types.EntitySoil && collision.Entity.Type == types.EntityPlayer {
		poistion := ecs.GetComponent[components.Position](collision.Entity)
		step := ecs.GetComponent[components.BlockStep](collision.Entity)

		block, err := game.NewEntityAndHeader(types.EntityVoid, world, s.stage, poistion.CurrentBlockPosition.Add(step.Increment))
		if err != nil {
			fmt.Printf("creating void block: %v\n", err)
			return
		}
		replaced := s.stage.VisitBlock(block, poistion.CurrentBlockPosition.Add(step.Increment))
		if replaced != nil && replaced != block {
			world.RemoveEntity(replaced)
		}
	}

	if collision.Block.Type == types.EntitySoil && collision.Entity.Type == types.EntityBoulder {
		s.sound.PlayFx("boulder_collision")
	}

	if collision.Block.Type == types.EntitySoil && collision.Entity.Type == types.EntityDiamond {
		s.sound.PlayFx("diamond_collision")
	}
}

func (s *GameplaySystem) handleEntityBlocks(world *ecs.World, event *game.EntityEvent) {
	// Magic wall: falling boulder/diamond hits magic wall entity → transform and pass through
	target := event.Target
	actor := event.Actor

	targetCharacteristic := ecs.GetComponent[components.Characteristic](target)
	if targetCharacteristic == nil || !targetCharacteristic.Has(common.IsMagicWall) {
		return
	}

	if actor.Type != types.EntityBoulder && actor.Type != types.EntityDiamond {
		return
	}

	// Check the entity was falling (direction is down)
	step := ecs.GetComponent[components.BlockStep](actor)
	if step == nil || step.Direction.Y <= 0 {
		return
	}

	// Determine output position (one cell below the magic wall entity)
	targetPosition := ecs.GetComponent[components.Position](target)
	if targetPosition == nil {
		return
	}
	outputPos := targetPosition.CurrentBlockPosition.Add(common.DirectionDown)

	// Check if the output position is empty
	if s.stage.CheckPositionOccupied(outputPos) {
		return
	}

	// Activate the magic wall timer and switch animations
	if !s.stage.IsMagicWallActive() {
		s.stage.ActivateMagicWall()
		s.setMagicWallAnimation(world, "active")
	}

	// Remove the falling entity
	actorPosition := ecs.GetComponent[components.Position](actor)
	s.stage.RemoveEntity(actor, actorPosition.CurrentBlockPosition)
	world.EnqueueRemoval(actor)

	// Spawn the transformed entity below the magic wall
	var newType ecs.EntityType
	if actor.Type == types.EntityBoulder {
		newType = types.EntityDiamond
	} else {
		newType = types.EntityBoulder
	}

	newEntity, err := game.NewEntityAndHeader(newType, world, s.stage, outputPos)
	if err == nil {
		s.stage.SetEntity(newEntity, outputPos)
		// Gravity will detect empty space below and start it falling naturally

		if newType == types.EntityDiamond {
			s.sound.PlayFx("diamond_created")
		}
	}

	fmt.Printf("Magic wall: %s transformed at %v\n", actor.ID, outputPos)
}

func (s *GameplaySystem) setMagicWallAnimation(world *ecs.World, animName string) {
	for _, entity := range world.Entities() {
		if entity.Type == types.EntityMagicWall {
			anim := ecs.GetComponent[components.Animation](entity)
			if anim != nil {
				anim.Current = animName
			}

			// When active, remove CanHoldGravity so objects above fall into the wall.
			// When inactive/expired, restore it so objects rest on top.
			characteristic := ecs.GetComponent[components.Characteristic](entity)
			if characteristic != nil {
				if animName == "active" {
					characteristic.Remove(common.CanHoldGravity)
				} else {
					characteristic.Add(common.CanHoldGravity)
				}
			}
		}
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

func (s *GameplaySystem) handleExitOpen(world *ecs.World) {

	if exit, ok := s.stage.GetEntityAtPosition(s.stage.ExitPosition); ok {
		ecs.RemoveComponent[components.Collider](exit)
		ecs.RemoveComponent[components.Sprite](exit)

		exit.AddComponent(components.NewSprite(common.NewSprite(assets.LoadTexture("entities"), s.stage.BlockWidth, s.stage.BlockHeight, float32(types.EntityExitDoor)*s.stage.BlockWidth, 0, 0)))
		exit.AddComponent(components.NewCollider(collider.LayerTrigger, collider.LayerNone, collider.LayerPlayer))

		door := ecs.GetComponent[components.Door](exit)
		door.State = components.DoorOpen
	}

	game.NewFlash(world)
	s.sound.PlayFx("stage_exit_opened")
}

func (s *GameplaySystem) handleDamage(world *ecs.World, damage *game.DamageEvent) {

	if damage.Target.Type == types.EntityPlayer {
		s.DamageOnPlayer(world, damage.Target)
	}

	if damage.Target.Type == types.EntityFirefly || damage.Target.Type == types.EntityButterfly {
		s.OnDamageOnEnemy(world, damage.Source, damage.Target)
	}
}

func (s *GameplaySystem) OnDamageOnEnemy(world *ecs.World, boulder *ecs.Entity, enemy *ecs.Entity) {
	fmt.Println("Boulder and enemy collision detected!!")

	boulderStep := ecs.GetComponent[components.BlockStep](boulder)

	if boulderStep.Direction.IsZero() {
		return
	}

	enemyPosition := ecs.GetComponent[components.Position](enemy)

	s.stage.TryRemoveEntity(enemy)
	world.EnqueueRemoval(enemy)

	s.stage.TryRemoveEntity(boulder)
	world.EnqueueRemoval(boulder)

	position := enemyPosition.CurrentBlockPosition

	if enemy.Type == types.EntityFirefly {
		s.sound.PlayFx("explosion")
		s.CreateSquare(world, types.EntityExplosion, position, 1, 1)
	}

	if enemy.Type == types.EntityButterfly {
		s.CreateSquare(world, types.EntityDiamond, position, 1, 1)
	}
}

func (s *GameplaySystem) OnBoulderPlayerCollision(world *ecs.World, boulder *ecs.Entity, enemy *ecs.Entity) {
	fmt.Println("Boulder and enemy collision detected!!")

	boulderStep := ecs.GetComponent[components.BlockStep](boulder)

	if boulderStep.Direction.IsZero() {
		return
	}

	enemyPosition := ecs.GetComponent[components.Position](enemy)

	s.stage.TryRemoveEntity(enemy)
	world.EnqueueRemoval(enemy)

	s.stage.TryRemoveEntity(boulder)
	world.EnqueueRemoval(boulder)

	position := enemyPosition.CurrentBlockPosition

	s.CreateDiamonds(world, position, 2, 2)
}

func (s *GameplaySystem) DamageOnPlayer(world *ecs.World, player *ecs.Entity) {
	fmt.Println("Player damage triggered!!")

	s.sound.PlayFx("player_hurt")
	s.sound.PlayFx("player_died")

	playerPosition := ecs.GetComponent[components.Position](player)
	position := playerPosition.CurrentBlockPosition

	world.EnqueueRemoval(player)
	s.CreateDiamonds(world, position, 2, 2)

	world.AddEvent("stagechange", game.NewStageChangeEvent(s.stage.Name, s.stage.EnterPosition))
}

func (s *GameplaySystem) CreateDiamonds(world *ecs.World, position common.BlockPosition, width int, height int) {
	for y := -height; y <= height; y++ {
		for x := -width; x <= width; x++ {
			diamondPosition := position.Offset(x, y)
			if s.stage.CheckCharacteristics(diamondPosition, common.Destructable) {
				block, err := game.NewEntityAndHeader(types.EntityVoid, world, s.stage, diamondPosition)
				if err != nil {
					fmt.Printf("creating void block at %d:%d: %v\n", diamondPosition.X, diamondPosition.Y, err)
					continue
				}
				old := s.stage.BlockMap.SetBlock(block, diamondPosition)
				if old != nil && old != block {
					world.RemoveEntity(old)
				}

				doomed := s.stage.GetEntity(diamondPosition)
				if doomed != nil {
					world.EnqueueRemoval(doomed)
				}

				diamond, err := game.NewEntityAndHeader(types.EntityDiamond, world, s.stage, diamondPosition)

				if err == nil {
					s.stage.SetEntity(diamond, diamondPosition)
				}
			}
		}
	}
}

func (s *GameplaySystem) CreateSquare(world *ecs.World, entityType ecs.EntityType, center common.BlockPosition, width int, height int) {
	for y := -height; y <= height; y++ {
		for x := -width; x <= width; x++ {
			pos := center.Offset(x, y)
			if s.stage.CheckCharacteristics(pos, common.Destructable) {
				block, err := game.NewEntityAndHeader(types.EntityVoid, world, s.stage, pos)
				if err != nil {
					fmt.Printf("creating void block at %d:%d: %v\n", pos.X, pos.Y, err)
					continue
				}
				old := s.stage.BlockMap.SetBlock(block, pos)
				if old != nil && old != block {
					world.RemoveEntity(old)
				}

				doomed := s.stage.GetEntity(pos)
				if doomed != nil {
					world.EnqueueRemoval(doomed)
				}

				entity, err := game.NewEntityAndHeader(entityType, world, s.stage, pos)

				if err == nil {
					s.stage.SetEntity(entity, pos)
				}
			}
		}
	}
}
