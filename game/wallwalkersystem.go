package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/core"
)

type wallWalkerWorld interface {
	CheckCharacteristics(core.BlockPosition, characteristics.Characteristics) bool
	GetPosition(core.BlockPosition) rl.Vector2
	Entities() []*Entity
}

type WallWalkerSystem struct {
	world wallWalkerWorld
}

func NewWallWalkerSystem(w *World) *WallWalkerSystem {
	return &WallWalkerSystem{
		world: w,
	}
}

func (s *WallWalkerSystem) Update() {
	for _, entity := range s.world.Entities() {
		if entity != nil && entity.Position != nil && entity.Velocity != nil && entity.WallWalker != nil {
			s.UpdateTarget(entity, entity.Position, entity.Velocity, entity.WallWalker)
		}
	}
}

func (s *WallWalkerSystem) UpdateTarget(entity *Entity, position *components.PositionComponent, velocity *components.VelocityComponent, walker *components.WallWalkerComponent) {
	current := position.CurrentBlockPosition
	w := s.world

	if !entity.HasCharacteristic(characteristics.IsEnemy) {
		return
	}

	if position.Vector2 != w.GetPosition(current) {
		return
	}

	if velocity.BlockVector.IsZero() {

		direction, _ := s.findDirection(current, core.BlockVector{X: 1, Y: 0})
		velocity.BlockVector = direction

		return

	} else if !position.HasTarget() {

		// Check surrounding positions relative to current direction
		rightPos := current.Add(velocity.BlockVector.TurnRight()) // Right side check
		aheadPos := current.Add(velocity.BlockVector)             // Forward check

		// Check for obstacles (non-void blocks or entities)
		hasRightObstacle := w.CheckCharacteristics(rightPos, characteristics.EnemyObstacle)
		hasAheadObstacle := w.CheckCharacteristics(aheadPos, characteristics.EnemyObstacle)

		if !hasRightObstacle && walker.HasWall {
			// Was following a wall but lost contact, turn right
			velocity.BlockVector = velocity.BlockVector.TurnRight()
		} else if hasAheadObstacle {
			// Wall ahead, turn left
			velocity.BlockVector = velocity.BlockVector.TurnLeft()
		}

		// else keep moving forward until we find a wall

		// Update wall following state
		entity.WallWalker.HasWall = hasRightObstacle
	}
}

func (s *WallWalkerSystem) findDirection(position core.BlockPosition, direction core.BlockVector) (core.BlockVector, bool) {
	for i := 0; i < 4; i++ {
		if !s.world.CheckCharacteristics(position.Add(direction), characteristics.EnemyObstacle) {
			return direction, true
		}

		direction = direction.TurnRight()
	}

	return direction, false
}
