package systems

import (
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

// type wallWalkerWorld interface {
// 	CheckCharacteristics(common.BlockPosition, characteristics.Characteristics) bool
// 	GetPosition(common.BlockPosition) rl.Vector2
// 	Entities() []*entity.Entity
// }

type WallWalkerSystem struct {
	stage *game.Stage
	PositionResolver
	BlockResolver
}

func NewWallWalkerSystem(stage *game.Stage) *WallWalkerSystem {
	return &WallWalkerSystem{stage: stage}
}

func (s *WallWalkerSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		comps := world.GetComponents(entity)
		position := ecs.GetComponent[components.PositionComponent](comps)
		step := ecs.GetComponent[components.BlockStep](comps)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)
		wallwalker := ecs.GetComponent[components.WallWalkerComponent](comps)

		if position != nil && step != nil && wallwalker != nil {
			s.UpdateTarget(world, position, step, characteristic, wallwalker)
		}
	}
}

func (s *WallWalkerSystem) UpdateTarget(world *ecs.World, position *components.PositionComponent, step *components.BlockStep, characteristic *components.CharacteristicComponent, wallwalker *components.WallWalkerComponent) {
	current := position.CurrentBlockPosition

	if !characteristic.Has(characteristics.IsEnemy) {
		return
	}

	if position.Vector2 != s.stage.GetPosition(current) {
		return
	}

	if step.Increment.IsZero() {

		direction, _ := s.findDirection(world, current, common.BlockVector{X: 1, Y: 0})
		step.Increment = direction

		return

	} else { // if !position.HasTarget()

		// Check surrounding positions relative to current direction
		rightPos := current.Add(step.Increment.TurnRight()) // Right side check
		aheadPos := current.Add(step.Increment)             // Forward check

		// Check for obstacles (non-void blocks or entities)
		hasRightObstacle := s.stage.CheckCharacteristics(world, rightPos, characteristics.EnemyObstacle)
		hasAheadObstacle := s.stage.CheckCharacteristics(world, aheadPos, characteristics.EnemyObstacle)

		if !hasRightObstacle && wallwalker.HasWall {
			// Was following a wall but lost contact, turn right
			step.Increment = step.Increment.TurnRight()
		} else if hasAheadObstacle {
			// Wall ahead, turn left
			step.Increment = step.Increment.TurnLeft()
		}

		// else keep moving forward until we find a wall

		// Update wall following state
		wallwalker.HasWall = hasRightObstacle
	}
}

func (s *WallWalkerSystem) findDirection(world *ecs.World, position common.BlockPosition, direction common.BlockVector) (common.BlockVector, bool) {
	for i := 0; i < 4; i++ {
		if !s.stage.CheckCharacteristics(world, position.Add(direction), characteristics.EnemyObstacle) {
			return direction, true
		}

		direction = direction.TurnRight()
	}

	return direction, false
}
