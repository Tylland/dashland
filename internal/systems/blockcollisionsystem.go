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
	// s.checkBlockCollisions()
	// if !s.Player().IsDead {
	// 	s.CheckPlayerCollisions(world)
	// }

	for _, entity := range world.Entities() {
		comps := world.GetComponents(entity)
		position := ecs.GetComponent[components.PositionComponent](comps)
		step := ecs.GetComponent[components.BlockStep](comps)
		collider := ecs.GetComponent[components.ColliderComponent](comps)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)

		if position != nil && step != nil && collider != nil {
			s.checkEntityCollisions(world, entity, position, step, collider, characteristic)
		}

	}
}

func (s *BlockCollisionSystem) stopFalling(entity *ecs.Entity, velocity *components.VelocityComponent, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	velocity.Clear()
	characteristic.Add(characteristics.PlayerObstacle)
	characteristic.Add(characteristics.EnemyObstacle)
}

func (s *BlockCollisionSystem) checkEntityCollisions(world *ecs.World, entity *ecs.Entity, position *components.PositionComponent, step *components.BlockStep, collider *components.ColliderComponent, characteristic *components.CharacteristicComponent) {

	// if position.CurrentBlockPosition.IsSame(s.Player().Position.CurrentBlockPosition) {
	// 	// Dispatch player collision event
	// 	s.stage.OnEvent(game.NewPlayerCollisionEvent(world, s.Player(), entity, position, velocity.IsFalling()))
	// 	return
	// }

	if step.Increment.IsZero() {
		return
	}

	targetPosition := position.CurrentBlockPosition.Add(step.Increment)

	if block, ok := s.GetBlockAtPosition(targetPosition); ok {
		if collider.CollidesWith(block.Collider) {
			world.AddComponent(entity, components.NewEvent("blockcollision", game.NewBlockCollisionEvent(world, block, entity)))
		}
	}

	if existing, ok := s.stage.GetEntityAtPosition(targetPosition); ok {
		existingComps := world.GetComponents(existing)
		existingCollider := ecs.GetComponent[components.ColliderComponent](existingComps)

		if existingCollider != nil && collider.CollidesWith(existingCollider) {
			world.AddComponent(entity, components.NewEvent("entitycollision", game.NewEntityCollisionEvent(world, entity, existing)))
		}
	}
}

// func (s *BlockCollisionSystem) CheckPlayerCollisions(world *ecs.World) {
// 	// if !s.Player().Position.HasTarget() {
// 	// 	return
// 	// }

// 	player := world.GetEntity("player")
// 	playerComps := world.GetComponents(player)
// 	playerPosition := ecs.GetComponent[components.PositionComponent](playerComps)
// 	playerVelocity := ecs.GetComponent[components.VelocityComponent](playerComps)

// 	playerTargetPos := playerPosition.GetBlockTarget(playerVelocity)

// 	if !s.CheckBlockAtPosition(game.Void, playerTargetPos) {
// 		return
// 	}

// 	entity := s.stage.GetEntityAtPosition(playerTargetPos)

// 	if entity == nil {
// 		return
// 	}

// 	comps := world.GetComponents(entity)

// 	characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)
// 	velocity := ecs.GetComponent[components.VelocityComponent](comps)

// 	if characteristic.Has(characteristics.Collectable) {
// 		// Only collect if the diamond is not falling
// 		if !velocity.IsFalling() {
// 			position := ecs.GetComponent[components.PositionComponent](comps)
// 			s.stage.OnEvent(game.NewPlayerCollisionEvent(world, player, entity, position, false))
// 		}
// 	}
//}
