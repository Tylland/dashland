package systems

import (
	"fmt"
	"math/rand"

	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type GravityBehavior struct {
	stage *game.Stage
}

func NewGravityBehavior(stage *game.Stage) *GravityBehavior {
	return &GravityBehavior{
		stage: stage,
	}
}

func (g *GravityBehavior) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		position := ecs.GetComponent[components.PositionComponent](entity)
		step := ecs.GetComponent[components.BlockStep](entity)
		collider := ecs.GetComponent[components.ColliderComponent](entity)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](entity)

		if position != nil && step != nil && collider != nil && characteristic != nil {
			g.ApplyGravityOnEntity(world, entity, position, step, collider, characteristic)
		}
	}
}

func (g *GravityBehavior) StartFalling(entity *ecs.Entity, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s start falling!!\n", entity.ID)

	characteristic.Remove(characteristics.Obstacle)
	characteristic.Add(characteristics.Falling)
}

func (g *GravityBehavior) StopFalling(entity *ecs.Entity, step *components.BlockStep, characteristic *components.CharacteristicComponent) {
	fmt.Printf("Entity %s stop falling!!\n", entity.ID)

	step.Increment.Clear()
	characteristic.Add(characteristics.Obstacle)
	characteristic.Remove(characteristics.Falling)
}

func (g *GravityBehavior) ApplyGravityOnEntity(world *ecs.World, entity *ecs.Entity, position *components.PositionComponent, step *components.BlockStep, collider *components.ColliderComponent, characteristic *components.CharacteristicComponent) {
	current := position.CurrentBlockPosition

	if !characteristic.Has(characteristics.CanFall) {
		return
	}

	// Don't apply gravity if entity is still moving to its target position
	if position.Vector2 != g.stage.GetPosition(current) {
		return
	}

	// If already falling
	if characteristic.Has(characteristics.Falling) {
		step.Move(common.DirectionDown, moveSpeed)
		fmt.Printf("%s continue falling", entity.ID)

		return
	}

	under := current.Add(common.DirectionDown)

	// Start falling if no support below
	if !g.stage.CheckCharacteristics(world, under, characteristics.CanHoldGravity) {
		step.Move(common.DirectionDown, moveSpeed)
		g.StartFalling(entity, characteristic)
		return
	}

	if !g.stage.CheckCharacteristics(world, under, characteristics.GravityRollOff) {
		return
	}

	var firstSide common.BlockPosition
	var firstSideUnder common.BlockPosition
	var secondSide common.BlockPosition
	var secondSideUnder common.BlockPosition

	if rand.Float64() > 0.5 {
		firstSide = current.Add(common.DirectionRight)
		firstSideUnder = current.Add(common.DirectionRightDown)

		secondSide = current.Add(common.DirectionLeft)
		secondSideUnder = current.Add(common.DirectionLeftDown)
	} else {
		firstSide = current.Add(common.DirectionLeft)
		firstSideUnder = current.Add(common.DirectionLeftDown)

		secondSide = current.Add(common.DirectionRight)
		secondSideUnder = current.Add(common.DirectionRightDown)
	}

	if !g.stage.CheckCharacteristics(world, firstSide, characteristics.CanHoldGravity) && !g.stage.CheckCharacteristics(world, firstSideUnder, characteristics.CanHoldGravity) {
		step.Move(common.DirectionRightDown, moveSpeed)
		return
	}

	if !g.stage.CheckCharacteristics(world, secondSide, characteristics.CanHoldGravity) && !g.stage.CheckCharacteristics(world, secondSideUnder, characteristics.CanHoldGravity) {
		step.Move(common.DirectionLeftDown, moveSpeed)
		return
	}
}
