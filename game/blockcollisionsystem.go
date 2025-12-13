package game

import "github.com/tylland/dashland/game/characteristics"

type BlockCollisionSystem struct {
	world *World
}

func NewBlockCollisionSystem(world *World) *BlockCollisionSystem {
	return &BlockCollisionSystem{
		world: world,
	}
}

func (s *BlockCollisionSystem) Update() {
	// s.checkBlockCollisions()
	if !s.world.player.IsDead {
		s.checkPlayerCollisions()
	}

	for _, entity := range s.world.entities {
		if entity != nil && entity.Position != nil && entity.Velocity != nil {
			s.checkEntityCollisions(entity)
		}

	}
}

func (s *BlockCollisionSystem) checkEntityCollisions(entity *Entity) {

	if entity.Position.CurrentBlockPosition.IsSame(s.world.player.Position.CurrentBlockPosition) {
		// Dispatch player collision event
		event := PlayerCollisionEvent{
			Player:        s.world.player,
			Entity:        entity,
			EntityFalling: entity.Velocity.IsFalling(),
		}

		s.world.OnEvent(event)
	}

	if entity.Velocity.BlockVector.IsZero() {
		return
	}

	position := entity.Position

	//Get block at entity position

	block, ok := s.world.GetBlockAtPosition(position.TargetBlockPosition)

	if ok && block.blockType != Void {
		entity.Velocity.BlockVector.Clear()
		entity.Position.CancelTarget()
		// TODO: Emit event for collision
		s.world.OnEvent(BlockCollisionEvent{Entity: entity, Block: block})
	}

	//Get entity at entity position
	entityAtPosition := s.world.GetEntityAtPosition(position.TargetBlockPosition)

	if entityAtPosition != nil && /* entityAtPosition.Type == Boulder && */ entity.Id != entityAtPosition.Id {
		entity.Velocity.BlockVector.Clear()
		entity.Position.CancelTarget()

		// TODO: Emit event for collision
		s.world.OnEvent(EntityCollisionEvent{Entity1: entity, Entity2: entityAtPosition})
	}

}

func (s *BlockCollisionSystem) checkPlayerCollisions() {
	if !s.world.player.Position.HasTarget() {
		return
	}

	targetPos := s.world.player.Position.TargetBlockPosition

	if !s.world.CheckBlockAtPosition(Void, targetPos) {
		return
	}

	entity := s.world.GetEntityAtPosition(targetPos)

	if entity != nil && entity.HasCharacteristic(characteristics.Collectable) {
		// Only collect if the diamond is not falling
		if !entity.Velocity.IsFalling() {
			event := PlayerCollisionEvent{
				Player:        s.world.player,
				Entity:        entity,
				EntityFalling: false,
			}
			s.world.OnEvent(event)
		}
	}
}
