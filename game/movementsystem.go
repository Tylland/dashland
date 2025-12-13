package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/components"
)

const moveSpeed = 128 // pixels per second

type BlockMovementSystem struct {
	world *World
}

func NewBlockMovementSystem(w *World) *BlockMovementSystem {
	return &BlockMovementSystem{
		world: w,
	}
}

func (s *BlockMovementSystem) Update(deltaTime float32) {

	if !s.world.player.IsDead {
		s.handlePlayerPush()
		s.updatePlayerMovement(deltaTime)
	}

	// Handle entities movement
	for _, entity := range s.world.entities {
		if entity != nil && entity.Position != nil && entity.Velocity != nil {
			s.updateBlockMovement(entity, deltaTime)
		}
	}

}

func (s *BlockMovementSystem) updatePlayerMovement(deltaTime float32) {
	p := s.world.player

	// if !p.Position.CurrentBlockPosition.IsSame(p.Position.TargetBlockPosition) {
	// 	p.Position.UseTarget())
	// }

	// s.handleBlockMovement(p.Position, p.Velocity, deltaTime)

	if !p.Position.CurrentBlockPosition.IsSame(p.Position.TargetBlockPosition) {
		p.movement = Movement{}

		targetPosition := p.game.world.GetPosition(p.Position.TargetBlockPosition)
		p.movement.Start(p.game.world.GetPosition(p.Position.CurrentBlockPosition), targetPosition, moveSpeed, nil)
		p.game.world.VisitBlock(p.Position.TargetBlockPosition)
		//		p.game.world.VisitObject(p, p.targetBlockPosition)

		p.Position.UseTarget()

		p.Position.Vector2 = targetPosition
	}

	if p.movement.moving {
		p.movement.Update(deltaTime)
		p.Position.Vector2 = p.movement.position
	}
}

func (s *BlockMovementSystem) handlePlayerPush() {
	// Get entity at player's target position

	if targetEntity := s.world.GetEntity(s.world.player.Position.TargetBlockPosition); targetEntity != nil {
		if !targetEntity.HasCharacteristic(characteristics.Pushable) {
			return
		}

		// Calculate push direction based on player's movement
		if s.world.player.Position.HasTarget() {
			pushPos := targetEntity.Position.TargetBlockPosition

			if s.world.player.Position.CurrentBlockPosition.X > targetEntity.Position.TargetBlockPosition.X {
				pushPos = pushPos.Offset(-1, 0)
			} else if s.world.player.Position.CurrentBlockPosition.X < targetEntity.Position.TargetBlockPosition.X {
				pushPos = pushPos.Offset(1, 0)
			}

			// Check if push position is free
			if s.world.CheckBlockAtPosition(Void, pushPos) && s.world.GetEntity(pushPos) == nil {
				// Move the pushed entity
				//s.world.MoveEntity(targetEntity, pushPos)
				targetEntity.Position.TargetBlockPosition = pushPos
			}
		}
	}
}

func (s *BlockMovementSystem) updateBlockMovement(entity *Entity, deltaTime float32) {
	// Regular movement code

	if !entity.Position.CurrentBlockPosition.IsSame(entity.Position.TargetBlockPosition) {
		s.world.MoveEntity(entity, entity.Position.TargetBlockPosition)
	}

	s.handleBlockMovement(entity.Position, entity.Velocity, deltaTime)
}

// Handle grid-based movement
func (s *BlockMovementSystem) handleBlockMovement(position *components.PositionComponent, velocity *components.VelocityComponent, deltaTime float32) {
	targetPos := s.world.GetPosition(position.CurrentBlockPosition)
	currentPos := position.Vector2

	if currentPos != targetPos {
		diff := rl.Vector2Subtract(targetPos, currentPos)
		length := rl.Vector2Length(diff)

		if length > 0 {
			moveAmount := moveSpeed * deltaTime

			if moveAmount >= length {
				position.Vector2 = targetPos
			} else {
				moveVec := rl.Vector2Scale(rl.Vector2Normalize(diff), moveAmount)
				position.Vector2 = rl.Vector2Add(currentPos, moveVec)
			}

		}
	} else if !velocity.BlockVector.IsZero() {
		position.SetTarget(velocity.BlockVector)
	}

}
