package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

// type GravityWorld interface {
// 	CheckBlockAtPosition(blockType game.BlockType, position common.BlockPosition) bool
// 	GetPosition(blockPosition common.BlockPosition) rl.Vector2
// 	checkPositionOccupied(position common.BlockPosition) bool
// 	checkPlayerAtPosition(position common.BlockPosition) bool
// }

// type GravityWorld interface {
// 	CheckBlockAtPosition(blockType BlockType, position common.BlockVector) bool
// 	GetPosition(blockPosition common.BlockPosition) common.Vector2
// }

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
		velocity := ecs.GetComponent[components.VelocityComponent](comps)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](comps)

		if position != nil && velocity != nil && characteristic != nil {
			g.ApplyGravityOnEntity(world, entity, position, velocity, characteristic)
		}
	}
}

func (g *GravitySystem) StartFalling(entity *ecs.Entity, velocity *components.VelocityComponent, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s start falling!!\n", entity.ID)

	characteristic.Remove(characteristics.PlayerObstacle)
	characteristic.Remove(characteristics.EnemyObstacle)
}

func (g *GravitySystem) StopFalling(entity *ecs.Entity, velocity *components.VelocityComponent, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	velocity.BlockVector.Clear()
	characteristic.Add(characteristics.PlayerObstacle)
	characteristic.Add(characteristics.EnemyObstacle)
}

func (g *GravitySystem) ApplyGravityOnEntity(world *ecs.World, entity *ecs.Entity, position *components.PositionComponent, velocity *components.VelocityComponent, characteristic *components.CharacteristicComponent) {
	current := position.CurrentBlockPosition

	if !characteristic.Has(characteristics.CanFall) {
		return
	}

	// Don't apply gravity if entity is still moving to its target position
	if position.Vector2 != g.stage.GetPosition(current) {
		return
	}

	// If already falling, continue with current velocity
	if velocity.BlockVector.Y > 0 {
		return
	}

	under := current.Offset(0, 1)

	// Start falling if no support below
	if !g.stage.CheckBlockAtPosition(game.Soil, under) && !g.stage.CheckCharacteristics(world, under, characteristics.CanHoldGravity) {
		if !g.stage.CheckPositionOccupied(under) {
			velocity.BlockVector.X = 0
			velocity.BlockVector.Y = 1
			g.StartFalling(entity, velocity, characteristic)
			return
		}

		// Try falling diagonally
		right := current.Offset(1, 0)
		rightUnder := current.Offset(1, 1)

		if !g.stage.CheckPositionOccupied(right) && !g.stage.CheckPositionOccupied(rightUnder) && !g.stage.CheckCharacteristics(world, rightUnder, characteristics.CanHoldGravity) {
			position.TargetBlockPosition = rightUnder
			g.StartFalling(entity, velocity, characteristic)
			return
		}

		left := current.Offset(-1, 0)
		leftUnder := current.Offset(-1, 1)

		if !g.stage.CheckPositionOccupied(left) && !g.stage.CheckPositionOccupied(leftUnder) && !g.stage.CheckCharacteristics(world, leftUnder, characteristics.CanHoldGravity) {
			position.TargetBlockPosition = leftUnder
			g.StartFalling(entity, velocity, characteristic)
			return
		}

	}

	// If we reach here, there's something blocking the fall
	if velocity.IsFalling() {
		g.StopFalling(entity, velocity, characteristic)
	}
}
