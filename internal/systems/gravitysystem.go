package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type GravitySystem struct {
	stage *game.Stage
	PositionResolver
	BlockResolver
}

func NewGravitySystem(stage *game.Stage) *GravitySystem {
	return &GravitySystem{
		stage: stage,
	}
}

func (g *GravitySystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		comps := world.GetComponents(entity)

		position := ecs.GetComponent[components.PositionComponent](comps)
		step := ecs.GetComponent[components.BlockStep](comps)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)

		if position != nil && step != nil && characteristic != nil {
			g.ApplyGravityOnEntity(world, entity, position, step, characteristic)
		}
	}
}

func (g *GravitySystem) StartFalling(entity *ecs.Entity, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s start falling!!\n", entity.ID)

	characteristic.Remove(characteristics.PlayerObstacle)
	characteristic.Remove(characteristics.EnemyObstacle)
}

func (g *GravitySystem) StopFalling(entity *ecs.Entity, step *components.BlockStep, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	step.Increment.Clear()
	characteristic.Add(characteristics.PlayerObstacle)
	characteristic.Add(characteristics.EnemyObstacle)
}

func (g *GravitySystem) ApplyGravityOnEntity(world *ecs.World, entity *ecs.Entity, position *components.PositionComponent, step *components.BlockStep, characteristic *components.CharacteristicComponent) {
	current := position.CurrentBlockPosition

	if !characteristic.Has(characteristics.CanFall) {
		return
	}

	// Don't apply gravity if entity is still moving to its target position
	if position.Vector2 != g.stage.GetPosition(current) {
		return
	}

	// If already falling, continue with current velocity
	if step.Increment.Y > 0 {
		return
	}

	under := current.Add(common.DirectionDown)

	// Start falling if no support below
	if !g.stage.CheckBlockAtPosition(game.Soil, under) && !g.stage.CheckCharacteristics(world, under, characteristics.CanHoldGravity) {
		if !g.stage.CheckPositionOccupied(under) {
			step.Move(common.DirectionDown, g.stage.GetPosition(under), moveSpeed)
			g.StartFalling(entity, characteristic)
			return
		}

		// Try falling diagonally
		right := current.Add(common.DirectionRight)
		rightUnder := current.Add(common.DirectionRightDown)

		if !g.stage.CheckPositionOccupied(right) && !g.stage.CheckPositionOccupied(rightUnder) && !g.stage.CheckCharacteristics(world, rightUnder, characteristics.CanHoldGravity) {
			step.Move(common.DirectionRightDown, g.stage.GetPosition(rightUnder), moveSpeed)
			g.StartFalling(entity, characteristic)
			return
		}

		left := current.Offset(-1, 0)
		leftUnder := current.Offset(-1, 1)

		if !g.stage.CheckPositionOccupied(left) && !g.stage.CheckPositionOccupied(leftUnder) && !g.stage.CheckCharacteristics(world, leftUnder, characteristics.CanHoldGravity) {
			step.Move(common.DirectionLeftDown, g.stage.GetPosition(leftUnder), moveSpeed)
			g.StartFalling(entity, characteristic)
			return
		}

	}

	// // If we reach here, there's something blocking the fall
	// if step.Increment() {
	// 	g.StopFalling(entity, step, characteristic)
	// }
}
