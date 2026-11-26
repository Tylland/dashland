package systems

import (
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
		position := ecs.GetComponent[components.Position](entity)
		step := ecs.GetComponent[components.BlockStep](entity)
		collider := ecs.GetComponent[components.Collider](entity)
		characteristic := ecs.GetComponent[components.Characteristic](entity)
		wallwalker := ecs.GetComponent[components.WallWalker](entity)

		if position != nil && step != nil && collider != nil && characteristic != nil && wallwalker != nil {
			s.UpdateTarget(world, position, step, collider, characteristic, wallwalker)
		}
	}
}

func (s *WallWalkerBehavior) UpdateTarget(world *ecs.World, position *components.Position, step *components.BlockStep, collider *components.Collider, characteristic *components.Characteristic, wallwalker *components.WallWalker) {
	current := position.CurrentBlockPosition

	if !characteristic.Has(common.IsEnemy) {
		return
	}

	if position.Vector2 != s.stage.GetPosition(current) {
		return
	}

	if step.Direction.IsZero() {
		direction := common.DirectionDown

		if wallwalker.Clockwise {
			direction = common.DirectionLeft
		}

		s.makeMove(step, direction)

		return
	} else if step.Increment.IsZero() {

		var sideTurn common.BlockVector

		if wallwalker.Clockwise {
			sideTurn = step.Direction.TurnLeft()
		} else {
			sideTurn = step.Direction.TurnRight()
		}

		sidePos := current.Add(sideTurn)
		aheadPos := current.Add(step.Direction)

		// Check for obstacles (non-void blocks or entities)
		hasSideObstacle := s.stage.CheckBlocked(world, sidePos, collider)
		hasAheadObstacle := s.stage.CheckBlocked(world, aheadPos, collider)

		if !hasSideObstacle && wallwalker.HasWall {
			// Was following a wall but lost contact, turn
			s.makeMove(step, sideTurn)
		} else if hasAheadObstacle {
			// Wall ahead, turn
			if wallwalker.Clockwise {
				sideTurn = step.Direction.TurnRight()
			} else {
				sideTurn = step.Direction.TurnLeft()
			}

			s.makeMove(step, sideTurn)
		}

		if step.Increment.IsZero() {
			s.makeMove(step, step.Direction)
		}

		// Update wall following state
		wallwalker.HasWall = hasSideObstacle
	}
}

func (s *WallWalkerBehavior) makeMove(step *components.BlockStep, move common.BlockVector) {
	step.Move(move, moveSpeed)
}
