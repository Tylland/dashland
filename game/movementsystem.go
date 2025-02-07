package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const moveSpeed = 128 // pixels per second

type BlockMovementSystem struct {
	world *world
}

func NewBlockMovementSystem(w *world) *BlockMovementSystem {
	return &BlockMovementSystem{
		world: w,
	}
}

func (s *BlockMovementSystem) Update(deltaTime float32) {
	// Handle entities movement
	for _, entity := range s.world.entities {
		if entity != nil && entity.Position != nil && entity.Velocity != nil {
			s.updateEntityMovement(entity, deltaTime)
		}
	}

	// Handle player pushing entities
	s.handlePlayerPush()
}

func (s *BlockMovementSystem) handlePlayerPush() {
	// Get entity at player's target position
	if targetEntity := s.world.GetEntity(s.world.player.Position.TargetBlockPosition); targetEntity != nil {
		if targetEntity.Behavior&Pushable == 0 {
			return
		}

		// Calculate push direction based on player's movement
		pushPos := targetEntity.Position.TargetBlockPosition

		if s.world.player.movement.moving {
			if s.world.player.Position.CurrentBlockPosition.X > targetEntity.Position.TargetBlockPosition.X {
				pushPos = pushPos.Offset(-1, 0)
			} else if s.world.player.Position.CurrentBlockPosition.X < targetEntity.Position.TargetBlockPosition.X {
				pushPos = pushPos.Offset(1, 0)
			}

			// Check if push position is free
			if s.world.CheckBlockAtPosition(Void, pushPos) && s.world.GetEntity(pushPos) == nil {
				// Move the pushed entity
				s.world.MoveEntity(targetEntity, pushPos)
			}
		}
	}
}

func (s *BlockMovementSystem) updateEntityMovement(entity *Entity, deltaTime float32) {
	// Regular movement code

	if !entity.Position.CurrentBlockPosition.IsSame(entity.Position.TargetBlockPosition) {
		s.world.MoveEntity(entity, entity.Position.TargetBlockPosition)
	}

	targetPos := s.world.GetPosition(entity.Position.CurrentBlockPosition)
	currentPos := entity.Position.Vector2

	if currentPos != targetPos {
		diff := rl.Vector2Subtract(targetPos, currentPos)
		length := rl.Vector2Length(diff)

		if length > 0 {
			moveAmount := moveSpeed * deltaTime

			if moveAmount >= length {
				entity.Position.Vector2 = targetPos
			} else {
				moveVec := rl.Vector2Scale(rl.Vector2Normalize(diff), moveAmount)
				entity.Position.Vector2 = rl.Vector2Add(currentPos, moveVec)
			}
		}
	} else if !entity.Velocity.BlockVector.IsZero() {

		// Handle grid-based movement
		entity.Position.SetTarget(entity.Velocity.BlockVector)

	}
}
