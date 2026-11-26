package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
)

type BlockCollisionSystem struct {
	stage *game.Stage
}

func NewBlockCollisionSystem(stage *game.Stage) *BlockCollisionSystem {
	return &BlockCollisionSystem{
		stage: stage,
	}
}

func (s *BlockCollisionSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		position := ecs.GetComponent[components.Position](entity)
		step := ecs.GetComponent[components.BlockStep](entity)
		collider := ecs.GetComponent[components.Collider](entity)
		character := ecs.GetComponent[components.Characteristic](entity)

		if position != nil && step != nil && collider != nil && character != nil {
			s.checkEntityCollisions(world, entity, position, step, collider, character)
		}
	}
}

func (s *BlockCollisionSystem) checkEntityCollisions(world *ecs.World, actor *ecs.Entity, position *components.Position, step *components.BlockStep, collider *components.Collider, character *components.Characteristic) {
	if step.Increment.IsZero() {
		return
	}

	targetPosition := position.CurrentBlockPosition.Add(step.Increment)

	if block, ok := s.stage.GetBlockAtPosition(targetPosition); ok {
		blockCollider := ecs.GetComponent[components.Collider](block)
		if blockCollider != nil {
			obstructs, interacts := collider.Collide(blockCollider)

			if obstructs {
				fmt.Printf("%s is obstructed by %s\n", actor.ID, block.ID)

				if character.Has(common.Falling) {
					s.StopFalling(actor, step, character)
				} else {
					step.Cancel()
				}

				world.AddEvent("blockblocks", game.NewBlockEvent(block, actor))
			}

			if interacts {
				fmt.Printf("%s overlap %s\n", actor.ID, block.ID)
				world.AddEvent("blockcollision", game.NewBlockCollisionEvent(block, actor))
			}
		}
	}

	if target, ok := s.stage.GetEntityAtPosition(targetPosition); ok {
		targetCollider := ecs.GetComponent[components.Collider](target)

		if targetCollider != nil {
			obstructs, interacts := collider.Collide(targetCollider)

			if obstructs {
				fmt.Printf("%s is obstructed by %s\n", actor.ID, target.ID)

				if character.Has(common.Falling) {
					s.StopFalling(actor, step, character)
				} else {
					step.Cancel()
				}

				world.AddEvent("entityobstructed", game.NewEntityEvent(actor, target))
			}

			if interacts {
				// If player collides with a door trigger -> stage change
				if actor.Type == types.EntityPlayer && target.Type == types.EntityExitDoor {
					door := ecs.GetComponent[components.Door](target)

					if door != nil && door.State == components.DoorOpen {
						world.AddEvent("stagechange", game.NewStageChangeEvent(door.Stage, door.Position))
						return
					}
				}

				if character.Has(common.Falling) {
					fmt.Printf("%s is falling\n", actor.ID)
				}

				actorDamage := ecs.GetComponent[components.Damage](actor)
				existingHealth := ecs.GetComponent[components.Health](target)

				if actorDamage != nil && existingHealth != nil {
					world.AddEvent("damage", game.NewDamageEvent(actor, target, actorDamage, existingHealth))
					return
				}

				collectable := ecs.GetComponent[components.Collectable](target)

				if collectable != nil {
					world.AddEvent("collect", game.NewCollectEvent(actor, target))
					return
				}

				world.AddEvent("entitycollision", game.NewEntityCollisionEvent(actor, target))
			}
		}
	}
}

func (s *BlockCollisionSystem) StopFalling(entity *ecs.Entity, step *components.BlockStep, characteristic *components.Characteristic) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	step.Cancel()
	characteristic.Remove(common.Falling)
}
