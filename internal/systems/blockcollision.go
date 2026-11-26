package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/characteristics"
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
		position := ecs.GetComponent[components.PositionComponent](entity)
		step := ecs.GetComponent[components.BlockStep](entity)
		collider := ecs.GetComponent[components.ColliderComponent](entity)
		character := ecs.GetComponent[components.CharacteristicComponent](entity)

		if position != nil && step != nil && collider != nil && character != nil {
			s.checkEntityCollisions(world, entity, position, step, collider, character)
		}
	}
}

func (s *BlockCollisionSystem) checkEntityCollisions(world *ecs.World, actor *ecs.Entity, position *components.PositionComponent, step *components.BlockStep, collider *components.ColliderComponent, character *components.CharacteristicComponent) {
	if step.Increment.IsZero() {
		return
	}

	targetPosition := position.CurrentBlockPosition.Add(step.Increment)

	if block, ok := s.stage.GetBlockAtPosition(targetPosition); ok {
		blockCollider := ecs.GetComponent[components.ColliderComponent](block)
		if blockCollider != nil {
			obstructs, interacts := collider.Collide(blockCollider)

			if obstructs {
				fmt.Printf("%s is obstructed by %s\n", actor.ID, block.ID)

				if character.Has(characteristics.Falling) {
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
		targetCollider := ecs.GetComponent[components.ColliderComponent](target)

		if targetCollider != nil {
			obstructs, interacts := collider.Collide(targetCollider)

			if obstructs {
				fmt.Printf("%s is obstructed by %s\n", actor.ID, target.ID)

				if character.Has(characteristics.Falling) {
					s.StopFalling(actor, step, character)
				} else {
					step.Cancel()
				}

				world.AddEvent("entityobstructed", game.NewEntityEvent(actor, target))
			}

			if interacts {
				// If player collides with a door trigger -> stage change
				if actor.Type == types.EntityPlayer && target.Type == types.EntityExitDoor {
					door := ecs.GetComponent[components.DoorComponent](target)

					if door != nil && door.State == components.DoorOpen {
						world.AddEvent("stagechange", game.NewStageChangeEvent(door.Stage, door.Position))
						return
					}
				}

				if character.Has(characteristics.Falling) {
					fmt.Printf("%s is falling\n", actor.ID)
				}

				actorDamage := ecs.GetComponent[components.Damage](actor)
				existingHealth := ecs.GetComponent[components.Health](target)

				if actorDamage != nil && existingHealth != nil {
					world.AddEvent("damage", game.NewDamageEvent(actor, target, actorDamage, existingHealth))
					return
				}

				collectable := ecs.GetComponent[components.CollectableComponent](target)

				if collectable != nil {
					world.AddEvent("collect", game.NewCollectEvent(actor, target))
					return
				}

				world.AddEvent("entitycollision", game.NewEntityCollisionEvent(actor, target))
			}
		}
	}
}

func (s *BlockCollisionSystem) StopFalling(entity *ecs.Entity, step *components.BlockStep, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	step.Cancel()
	characteristic.Remove(characteristics.Falling)
}
