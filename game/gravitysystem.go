package game

import (
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/components"
)

type GravitySystem struct {
	world *World
}

func NewGravitySystem(w *World) *GravitySystem {
	return &GravitySystem{
		world: w,
	}
}

func (g *GravitySystem) Update() {
	for _, entity := range g.world.entities {
		if entity != nil && entity.Position != nil && entity.Velocity != nil {
			g.ApplyGravityOnEntity(entity, entity.Position, entity.Velocity)
		}
	}
}

func (g *GravitySystem) ApplyGravityOnEntity(entity *Entity, position *components.PositionComponent, velocity *components.VelocityComponent) {
	current := position.CurrentBlockPosition
	w := g.world

	if !entity.HasCharacteristic(characteristics.CanFall) {
		return
	}

	// Don't apply gravity if entity is still moving to its target position
	if position.Vector2 != w.GetPosition(current) {
		return
	}

	// If already falling, continue with current velocity
	if velocity.BlockVector.Y > 0 {
		return
	}

	under := current.Offset(0, 1)

	// Start falling if no support below
	if !w.CheckBlockAtPosition(Soil, under) && !w.checkPlayerAtPosition(under) {
		if !w.checkPositionOccupied(under) {
			velocity.BlockVector.X = 0
			velocity.BlockVector.Y = 1
			return
		}

		// Try falling diagonally
		right := current.Offset(1, 0)
		rightUnder := current.Offset(1, 1)

		if !w.checkPositionOccupied(right) && !w.checkPositionOccupied(rightUnder) && !w.checkPlayerAtPosition(rightUnder) {
			entity.Position.TargetBlockPosition = rightUnder
			return
		}

		left := current.Offset(-1, 0)
		leftUnder := current.Offset(-1, 1)

		if !w.checkPositionOccupied(left) && !w.checkPositionOccupied(leftUnder) && !w.checkPlayerAtPosition(leftUnder) {
			entity.Position.TargetBlockPosition = leftUnder
			return
		}

	}

	// If we reach here, there's something blocking the fall
	velocity.BlockVector.Clear()
}
