package systems

import (
	"fmt"

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
		position := ecs.GetComponent[components.PositionComponent](entity.Components)
		step := ecs.GetComponent[components.BlockStep](entity.Components)
		collider := ecs.GetComponent[components.ColliderComponent](entity.Components)
		characteristic := ecs.GetComponent[components.CharacteristicComponent](entity.Components)

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

	// Try falling diagonally
	right := current.Add(common.DirectionRight)
	rightUnder := current.Add(common.DirectionRightDown)

	if !g.stage.CheckCharacteristics(world, right, characteristics.CanHoldGravity) && !g.stage.CheckCharacteristics(world, rightUnder, characteristics.CanHoldGravity) {
		step.Move(common.DirectionRightDown, moveSpeed)
		//g.StartFalling(entity, characteristic)
		return
	}

	left := current.Offset(-1, 0)
	leftUnder := current.Offset(-1, 1)

	if !g.stage.CheckCharacteristics(world, left, characteristics.CanHoldGravity) && !g.stage.CheckCharacteristics(world, leftUnder, characteristics.CanHoldGravity) {
		step.Move(common.DirectionLeftDown, moveSpeed)
		//g.StartFalling(entity, characteristic)
		return
	}
}
