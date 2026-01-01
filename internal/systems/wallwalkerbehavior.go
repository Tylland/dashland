package systems

import (
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type WallWalkerBehavior struct {
	stage *game.Stage
}

func NewWallWalkerBehavior(stage *game.Stage) *WallWalkerBehavior {
	return &WallWalkerBehavior{stage: stage}
}

func (s *WallWalkerBehavior) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		position := ecs.GetComponent[components.PositionComponent](entity)
		step := ecs.GetComponent[components.BlockStep](entity)
		collider := ecs.GetComponent[components.ColliderComponent](entity)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](entity)
		wallwalker := ecs.GetComponent[components.WallWalkerComponent](entity)

		if position != nil && step != nil && collider != nil && characteristic != nil && wallwalker != nil {
			s.UpdateTarget(world, position, step, collider, characteristic, wallwalker)
		}
	}
}

func (s *WallWalkerBehavior) UpdateTarget(world *ecs.World, position *components.PositionComponent, step *components.BlockStep, collider *components.ColliderComponent, characteristic *components.CharacteristicComponent, wallwalker *components.WallWalkerComponent) {
	current := position.CurrentBlockPosition

	if !characteristic.Has(characteristics.IsEnemy) {
		return
	}

	if position.Vector2 != s.stage.GetPosition(current) {
		return
	}

	if step.Direction.IsZero() {
		direction, _ := s.findDirection(world, current, common.DirectionDown)

		s.makeMove(step, direction)

		return
	} else if step.Increment.IsZero() {
		// Check surrounding positions relative to current direction
		rightTurn := step.Direction.TurnRight()

		rightPos := current.Add(rightTurn)      // Right side check
		aheadPos := current.Add(step.Direction) // Forward check

		// Check for obstacles (non-void blocks or entities)
		hasRightObstacle := s.stage.CheckBlocked(world, rightPos, collider)
		hasAheadObstacle := s.stage.CheckBlocked(world, aheadPos, collider)

		if !hasRightObstacle && wallwalker.HasWall {
			// Was following a wall but lost contact, turn right
			s.makeMove(step, rightTurn)
		} else if hasAheadObstacle {
			// Wall ahead, turn left
			leftTurn := step.Direction.TurnLeft()

			s.makeMove(step, leftTurn)
		}

		if step.Increment.IsZero() {
			s.makeMove(step, step.Direction)
		}

		// Update wall following state
		wallwalker.HasWall = hasRightObstacle
	}
}

func (s *WallWalkerBehavior) makeMove(step *components.BlockStep, move common.BlockVector) {
	step.Move(move, moveSpeed)
}

func (s *WallWalkerBehavior) findDirection(world *ecs.World, position common.BlockPosition, direction common.BlockVector) (common.BlockVector, bool) {
	for i := 0; i < 4; i++ {
		if !s.stage.CheckCharacteristics(position.Add(direction), characteristics.Obstacle) {
			return direction, true
		}

		direction = direction.TurnRight()
	}

	return direction, false
}
