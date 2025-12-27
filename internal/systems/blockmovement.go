package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

const moveSpeed = 128 // pixels per second

type BlockMovement struct {
	stage *game.Stage
}

func NewBlockMovement(stage *game.Stage) *BlockMovement {
	return &BlockMovement{
		stage: stage,
	}
}

func (s *BlockMovement) Update(world *ecs.World, deltaTime float32) {
	// Handle entities movement
	for _, entity := range world.Entities() {
		position := ecs.GetComponent[components.PositionComponent](entity)
		step := ecs.GetComponent[components.BlockStep](entity)

		if position != nil && step != nil {
			s.moveEntity(entity, position, step, deltaTime)
		}
	}
}

func (s *BlockMovement) moveEntity(entity *ecs.Entity, position *components.PositionComponent, step *components.BlockStep, deltaTime float32) {

	if !step.Increment.IsZero() {
		// if entity.ID == "player" {
		// 	fmt.Print("Player is here!!")
		// }

		target := position.CurrentBlockPosition.Add(step.Increment)

		if s.stage.TryMoveEntity(entity, position.CurrentBlockPosition, target) {
			position.Update(target)
			//			s.stage.VisitBlock(position.CurrentBlockPosition)
			step.Commit(s.stage.GetPosition(target))
		} else {
			character := ecs.GetComponent[components.CharacteristicComponent](entity)

			if character != nil {
				character.Remove(characteristics.Falling)
			}

			step.Cancel()
		}

	}

	s.moveToTarget(position, step, deltaTime)
}

func (s *BlockMovement) moveToTarget(position *components.PositionComponent, step *components.BlockStep, deltaTime float32) {
	targetPos := step.Target
	currentPos := position.Vector2

	if currentPos != targetPos {
		diff := rl.Vector2Subtract(targetPos, currentPos)
		length := rl.Vector2Length(diff)

		if length > 0 {
			moveAmount := step.Speed * deltaTime

			if moveAmount >= length {
				position.Vector2 = targetPos
			} else {
				moveVec := rl.Vector2Scale(rl.Vector2Normalize(diff), moveAmount)
				position.Vector2 = rl.Vector2Add(currentPos, moveVec)
			}
		}
	}
}
