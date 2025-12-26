package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type BlockCollisionSystem struct {
	stage *game.Stage
	BlockResolver
}

func NewBlockCollisionSystem(stage *game.Stage) *BlockCollisionSystem {
	return &BlockCollisionSystem{
		BlockResolver: stage,
		stage:         stage,
	}
}

func (s *BlockCollisionSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		position := ecs.GetComponent[components.PositionComponent](entity.Components)
		step := ecs.GetComponent[components.BlockStep](entity.Components)
		collider := ecs.GetComponent[components.ColliderComponent](entity.Components)
		character := ecs.GetComponent[components.CharacteristicComponent](entity.Components)

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

	if block, ok := s.GetBlockAtPosition(targetPosition); ok {
		blocks, collides := collider.Result(block.Collider)

		if blocks {
			fmt.Printf("%s is blocked by %s\n", actor.ID, block.ID)

			if character.Has(characteristics.Falling) {
				s.StopFalling(actor, step, character)
			} else {
				step.Cancel()
			}

			world.AddEvent("blockblocks", game.NewBlockEvent(block, actor))
		}

		if collides {
			fmt.Printf("%s overlap %s\n", actor.ID, block.ID)
			world.AddEvent("blockcollision", game.NewBlockCollisionEvent(block, actor))
		}
	}

	if target, ok := s.stage.GetEntityAtPosition(targetPosition); ok {
		targetCollider := ecs.GetComponent[components.ColliderComponent](target.Components)

		if targetCollider != nil {
			blocked, collides := collider.Result(targetCollider)

			if blocked {
				fmt.Printf("%s is blocked by %s\n", actor.ID, target.ID)

				if character.Has(characteristics.Falling) {
					s.StopFalling(actor, step, character)
				} else {
					step.Cancel()
				}

				world.AddEvent("entityblocks", game.NewEntityEvent(actor, target))
			}

			if collides {

				if character.Has(characteristics.Falling) {
					fmt.Printf("%s is falling\n", actor.ID)
				}

				// targetCharacter := ecs.GetComponent[components.CharacteristicComponent](target.Components)

				// if

				actorDamage := ecs.GetComponent[components.Damage](actor.Components)
				existingHealth := ecs.GetComponent[components.Health](target.Components)

				if actorDamage != nil && existingHealth != nil {
					world.AddEvent("damage", game.NewDamageEvent(actor, target, actorDamage, existingHealth))
					return
				}

				collectable := ecs.GetComponent[components.CollectableComponent](target.Components)

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
	characteristic.Add(characteristics.Obstacle)
	characteristic.Remove(characteristics.Falling)
}
