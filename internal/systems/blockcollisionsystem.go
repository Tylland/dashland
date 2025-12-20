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
		velocity := ecs.GetComponent[components.VelocityComponent](comps)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)

		if position != nil && velocity != nil {
			s.checkEntityCollisions(world, entity, position, velocity, characteristic)
		}

	}
}

func (s *BlockCollisionSystem) stopFalling(entity *ecs.Entity, velocity *components.VelocityComponent, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	velocity.BlockVector.Clear()
	characteristic.Add(characteristics.PlayerObstacle)
	characteristic.Add(characteristics.EnemyObstacle)
}

func (s *BlockCollisionSystem) checkEntityCollisions(world *ecs.World, entity *ecs.Entity, position *components.PositionComponent, velocity *components.VelocityComponent, characteristic *components.CharacteristicComponent) {

	// if position.CurrentBlockPosition.IsSame(s.Player().Position.CurrentBlockPosition) {
	// 	// Dispatch player collision event
	// 	s.stage.OnEvent(game.NewPlayerCollisionEvent(world, s.Player(), entity, position, velocity.IsFalling()))
	// 	return
	// }

	if velocity.BlockVector.IsZero() {
		return
	}

	//Get block at entity position

	block, ok := s.GetBlockAtPosition(position.TargetBlockPosition)

	if ok && block.BlockType != game.Void {
		s.stopFalling(entity, velocity, characteristic)

		position.CancelTarget()
		// TODO: Emit event for collision
		s.stage.OnEvent(game.NewBlockCollisionEvent(world, block, entity))
	}

	//Get entity at entity position
	entityAtPosition := s.stage.GetEntityAtPosition(position.TargetBlockPosition)

	if entityAtPosition != nil && /* entityAtPosition.Type == Boulder && */ entity.ID != entityAtPosition.ID {
		velocity.BlockVector.Clear()
		position.CancelTarget()

		// TODO: Emit event for collision
		s.stage.OnEvent(game.NewEntityCollisionEvent(world, entity, entityAtPosition))
	}

}

func (s *BlockCollisionSystem) CheckPlayerCollisions(world *ecs.World) {
	// if !s.Player().Position.HasTarget() {
	// 	return
	// }

	player := world.GetEntity("player")
	playerComps := world.GetComponents(player)
	position := ecs.GetComponent[components.PositionComponent](playerComps)

	targetPos := position.TargetBlockPosition

	if !s.CheckBlockAtPosition(game.Void, targetPos) {
		return
	}

	entity := s.stage.GetEntityAtPosition(targetPos)

	if entity == nil {
		return
	}

	comps := world.GetComponents(entity)

	characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)
	velocity := ecs.GetComponent[components.VelocityComponent](comps)

	if characteristic.Has(characteristics.Collectable) {
		// Only collect if the diamond is not falling
		if !velocity.IsFalling() {
			position := ecs.GetComponent[components.PositionComponent](comps)
			s.stage.OnEvent(game.NewPlayerCollisionEvent(world, player, entity, position, false))
		}
	}
}
